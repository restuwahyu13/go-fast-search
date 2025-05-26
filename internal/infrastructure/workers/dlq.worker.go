package worker

import (
	"context"
	"fmt"
	"strings"

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

type deadLetterQueueWorker struct {
	ctx  context.Context
	env  dto.Request[dto.Environtment]
	db   *bun.DB
	rds  *redis.Client
	amqp *rabbitmq.Conn
	mls  meilisearch.ServiceManager
}

func NewDeadLetterQueueWorker(options dto.WorkerOptions) inf.IDeadLetterQueueWorker {
	return deadLetterQueueWorker{ctx: options.CTX, env: options.ENV, db: options.DB, rds: options.RDS, amqp: options.AMQP, mls: options.MLS}
}

func (w deadLetterQueueWorker) deadLetterQueueRabbitInstance() inf.IRabbitMQ {
	return pkg.NewRabbitMQ(w.ctx, w.amqp)
}

func (w deadLetterQueueWorker) deadLetterQueueConsumer() {
	amqp := w.deadLetterQueueRabbitInstance()
	amqp_req := dto.Request[dto.RabbitOptions]{}

	amqp_req.Option.ExchangeName = cons.EXCHANGE_NAME_DEAD_LETTER_QUEUE
	amqp_req.Option.ExchangeType = cons.EXCHANGE_TYPE_DIRECT
	amqp_req.Option.QueueName = cons.QUEUE_NAME_DEAD_LETTER_QUEUE
	amqp_req.Option.Prefetch = 1

	amqp.Consumer(amqp_req, func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		if d.Headers[cons.X_RABBIT_UNKNOWN] == cons.TRUE && d.Headers[cons.X_RABBIT_SECRET] != w.env.Config.RABBITMQ.SECRET {
			pkg.Logrus(cons.INFO, "Queue is not allowed to be consumed: %s", string(d.Body))
			return rabbitmq.NackDiscard
		}

		parser := helper.NewParser()
		if err := parser.Unmarshal(d.Body, &amqp_req); err != nil {
			return rabbitmq.NackDiscard
		}

		if amqp_req.Option.Body != nil {
			amqp_req.Option.Args = rabbitmq.Table{cons.X_RABBIT_SECRET: w.env.Config.RABBITMQ.SECRET, cons.X_MESSAGE_TTL: 15}

			key := fmt.Sprintf("%s", strings.ToUpper(amqp_req.Option.QueueName))
			req := dto.Request[dto.SleepBackoff]{}

			req.Body.Ctx = w.ctx
			req.Config.Redis = w.rds
			req.Body.Key = key
			req.Body.Count = 1
			req.Body.Retry = 5
			req.Body.BackOffTime = 300

			_, err := helper.SleepBackoff[int](req, func() (int, error) {
				if err := amqp.Publisher(amqp_req); err != nil {
					return -1, err
				}

				pkg.Logrus(cons.INFO, "Queue is allowed to be consumed: %s", string(d.Body))
				return 1, nil
			})

			if err != nil {
				return rabbitmq.NackDiscard
			}

			return rabbitmq.Ack
		}

		return rabbitmq.NackDiscard
	})
}

func (w deadLetterQueueWorker) DeadLetterQueueRun() {
	w.deadLetterQueueConsumer()
}
