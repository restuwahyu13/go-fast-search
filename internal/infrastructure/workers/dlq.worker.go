package worker

import (
	"context"
	"math"
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

type workerDeadLetterQueue struct {
	ctx  context.Context
	env  dto.Request[dto.Environtment]
	db   *bun.DB
	rds  *redis.Client
	amqp *rabbitmq.Conn
	mls  meilisearch.ServiceManager
}

func NewDeadLetterQueueWorker(options dto.WorkerOptions) inf.IDeadLetterQueueWorker {
	return workerDeadLetterQueue{ctx: options.CTX, env: options.ENV, db: options.DB, rds: options.RDS, amqp: options.AMQP, mls: options.MLS}
}

func (w workerDeadLetterQueue) deadLetterQueueRabbitInstance() inf.IRabbitMQ {
	return pkg.NewRabbitMQ(w.ctx, w.amqp)
}

func (w workerDeadLetterQueue) deadLetterQueueConsumer() {
	amqp := w.deadLetterQueueRabbitInstance()
	amqp_req := dto.Request[dto.RabbitOptions]{}

	amqp_req.Option.ExchangeName = cons.EXCHANGE_NAME_DEAD_LETTER_QUEUE
	amqp_req.Option.ExchangeType = cons.EXCHANGE_TYPE_DIRECT
	amqp_req.Option.QueueName = cons.QUEUE_NAME_DEAD_LETTER_QUEUE
	amqp_req.Option.Prefetch = 1
	amqp_req.Option.Args = rabbitmq.Table{cons.X_RABBIT_SECRET: w.env.Config.RABBITMQ.SECRET}

	amqp.Consumer(amqp_req, func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		if d.Headers[cons.X_RABBIT_UNKNOWN] == cons.TRUE && d.Headers[cons.X_RABBIT_SECRET] != w.env.Config.RABBITMQ.SECRET {
			pkg.Logrus(cons.INFO, "Queue is not allowed to be consumed: %s", string(d.Body))
			return rabbitmq.NackDiscard
		}

		pkg.Logrus(cons.INFO, "Before queue is allowed to be consumed: %s", string(d.Body))

		parser := helper.NewParser()
		if err := parser.Unmarshal(d.Body, &amqp_req); err != nil {
			return rabbitmq.NackDiscard
		}

		count := 0
		retry := 10
		backoff := time.Duration(math.Pow(2, float64(count))) * time.Second

		if count <= retry {
			count++
			time.Sleep(backoff)
		}

		pkg.Logrus(cons.INFO, "After queue is allowed to be consumed: %#v", amqp_req.Option)

		if amqp_req.Option.Body != nil {
			if err := amqp.Publisher(amqp_req); err != nil {
				pkg.Logrus(cons.ERROR, err)
				return rabbitmq.NackDiscard
			}
		}

		return rabbitmq.Ack
	})
}

func (w workerDeadLetterQueue) signalDeadLetterQueue() inf.IRedis {
	now := time.Now().Format(cons.DATE_TIME_FORMAT)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM, syscall.SIGABRT, syscall.SIGUSR1)

	for {
		select {
		case sig := <-ch:
			pkg.Logrus(cons.INFO, "%s - Worker dql is received signal: %s", now, sig.String())

			if w.env.Config.APP.ENV != cons.DEV {
				time.Sleep(time.Second * 3)
			} else {
				time.Sleep(time.Second * 10)
			}

			os.Exit(0)

		default:
			time.Sleep(time.Second * 3)
			pkg.Logrus(cons.INFO, "%s - Worker dlq is running...", now)
		}
	}
}

func (w workerDeadLetterQueue) DeadLetterQueueRun() {
	w.deadLetterQueueConsumer()
	w.signalDeadLetterQueue()
}
