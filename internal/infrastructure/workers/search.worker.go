package worker

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/wagslane/go-rabbitmq"

	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type searchWorker struct {
	ctx  context.Context
	env  dto.Request[dto.Environtment]
	db   *bun.DB
	rds  *redis.Client
	amqp *rabbitmq.Conn
	mls  meilisearch.ServiceManager
}

func NewSearchWorker(options dto.WorkerOptions) inf.ISearchWorker {
	return searchWorker{ctx: options.CTX, env: options.ENV, db: options.DB, rds: options.RDS, amqp: options.AMQP, mls: options.MLS}
}

func (w searchWorker) searchRabbitInstance() inf.IRabbitMQ {
	return pkg.NewRabbitMQ(w.ctx, w.amqp)
}

func (w searchWorker) searchChangeDataCapture() error {
	key := "WORKER:SEARCH:CDC"
	value := time.Now().Format(cons.DATE_TIME_FORMAT)

	rds, err := pkg.NewRedis(w.ctx, w.rds)
	if err != nil {
		return err
	}

	isExist, err := rds.Exists(key)
	if err != nil {
		return err
	}

	if isExist < 1 {
		if err := rds.Set(key, value); err != nil {
			return err
		}
	}

	return nil
}

func (w searchWorker) searchDeadLetterQueue(amqp inf.IRabbitMQ, req *dto.RabbitDeadLetterQueueOptions) error {
	amqp_req := dto.Request[dto.RabbitOptions]{}
	amqp_req.Option.ExchangeName = req.Exchange
	amqp_req.Option.ExchangeType = req.ExchangeType
	amqp_req.Option.QueueName = req.Queue
	amqp_req.Option.Body = req.Body

	amqp_body := amqp_req

	amqp_req.Option.ExchangeName = cons.EXCHANGE_NAME_DEAD_LETTER_QUEUE
	amqp_req.Option.ExchangeType = cons.EXCHANGE_TYPE_DIRECT
	amqp_req.Option.QueueName = cons.QUEUE_NAME_DEAD_LETTER_QUEUE
	amqp_req.Option.Body = amqp_body
	amqp_req.Option.Args = rabbitmq.Table{
		cons.X_RABBIT_SECRET:  req.Secret,
		cons.X_RABBIT_UNKNOWN: req.Unknown,
		cons.X_MESSAGE_TTL:    15,
	}

	if err := amqp.Publisher(amqp_req); err != nil {
		return err
	}

	return nil
}

func (w searchWorker) searchHandler(req dto.Request[dto.MeiliSearchDocuments[map[string]any]]) error {
	mls := pkg.NewMeiliSearch(w.ctx, w.mls)

	switch req.Body.Action {

	case cons.INSERT:
		if _, err := mls.Insert(req.Body.Doc, &req.Body.Data); err != nil {
			return err
		}
		return nil

	case cons.UPDATE:
		if _, err := mls.Update(req.Body.Doc, req.Body.ID.(string), &req.Body.Data); err != nil {
			return err
		}
		return nil

	case cons.DELETE:
		if _, err := mls.Delete(req.Body.Doc, req.Body.ID.(string)); err != nil {
			return err
		}
		return nil

	default:
		return errors.New("Meilisearch unknown action")
	}
}

func (w searchWorker) searchConsumer() {
	amqp := w.searchRabbitInstance()
	amqp_req := dto.Request[dto.RabbitOptions]{}

	amqp_req.Option.ExchangeName = cons.EXCHANGE_NAME_SEARCH
	amqp_req.Option.ExchangeType = cons.EXCHANGE_TYPE_DIRECT
	amqp_req.Option.QueueName = cons.QUEUE_NAME_SEARCH
	amqp_req.Option.Prefetch = 1
	amqp_req.Option.Args = rabbitmq.Table{cons.X_RABBIT_SECRET: w.env.Config.RABBITMQ.SECRET}

	amqp.Consumer(amqp_req, func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		parser := helper.NewParser()

		dlq_req := dto.RabbitDeadLetterQueueOptions{}
		req := dto.Request[dto.MeiliSearchDocuments[map[string]any]]{}

		if err := parser.Unmarshal(d.Body, &req.Body); err != nil {
			return rabbitmq.NackDiscard
		}

		dlq_req.Body = req.Body

		if d.Headers[cons.X_RABBIT_SECRET] != w.env.Config.RABBITMQ.SECRET {
			dlq_req.Secret = cons.EMPTY
			dlq_req.Unknown = cons.TRUE
			dlq_req.Error = errors.New("Queue is not allowed to be consumed")
		}

		if err := w.searchChangeDataCapture(); err != nil {
			dlq_req.Secret = d.Headers[cons.X_RABBIT_SECRET]
			dlq_req.Unknown = cons.FALSE
			dlq_req.Error = err
		}

		if err := w.searchHandler(req); err != nil {
			dlq_req.Secret = d.Headers[cons.X_RABBIT_SECRET]
			dlq_req.Unknown = cons.FALSE
			dlq_req.Error = err
		}

		if dlq_req.Body.Data != nil && dlq_req.Error != nil {
			pkg.Logrus(cons.ERROR, dlq_req.Error)

			dlq_req.Exchange = amqp_req.Option.ExchangeName
			dlq_req.ExchangeType = amqp_req.Option.ExchangeType
			dlq_req.Queue = amqp_req.Option.QueueName

			if err := w.searchDeadLetterQueue(amqp, &dlq_req); err != nil {
				pkg.Logrus(cons.ERROR, err)
				return rabbitmq.NackDiscard
			}

			return rabbitmq.NackDiscard
		}

		return rabbitmq.Ack
	})
}

func (w searchWorker) searchSignal() {
	now := time.Now().Format(cons.DATE_TIME_FORMAT)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM, syscall.SIGABRT, syscall.SIGUSR1)

	for {
		select {
		case sig := <-ch:
			pkg.Logrus(cons.INFO, "%s - Worker search is received signal: %s", now, sig.String())

			if w.env.Config.APP.ENV != cons.DEV {
				time.Sleep(time.Second * 10)
			} else {
				time.Sleep(time.Second * 15)
			}

			os.Exit(0)

		default:
			time.Sleep(time.Second * 5)
			pkg.Logrus(cons.INFO, "%s - Worker search is running...", now)
		}
	}
}

func (w searchWorker) SearchRun() {
	w.searchConsumer()
	w.searchSignal()
}
