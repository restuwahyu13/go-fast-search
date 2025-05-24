package pkg

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/lithammer/shortuuid"
	amqp "github.com/wagslane/go-rabbitmq"

	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type rabbitmq struct {
	ctx      context.Context
	rabbitmq *amqp.Conn
}

func NewRabbitMQ(ctx context.Context, con *amqp.Conn) inf.IRabbitMQ {
	return rabbitmq{ctx: ctx, rabbitmq: con}
}

func (p rabbitmq) Publisher(req dto.Request[dto.RabbitOptions]) error {
	parser := helper.NewParser()

	if req.Option.ContentType == "" {
		req.Option.ContentType = "application/json"
	}

	if req.Option.Timestamp.Sub(time.Now()).Seconds() < 1 {
		req.Option.Timestamp = time.Now().Local()
	}

	publisher, err := amqp.NewPublisher(p.rabbitmq,
		amqp.WithPublisherOptionsExchangeName(req.Option.ExchangeName),
		amqp.WithPublisherOptionsExchangeKind(req.Option.ExchangeType),
		amqp.WithPublisherOptionsExchangeDeclare,
		amqp.WithPublisherOptionsExchangeDurable,
		amqp.WithPublisherOptionsExchangeNoWait,
		amqp.WithPublisherOptionsExchangeArgs(req.Option.Args),
		amqp.WithPublisherOptionsLogging,
	)

	defer p.closeConnection(publisher, nil, p.rabbitmq)
	if err != nil {
		return err
	}

	bodyByte, err := parser.Marshal(&req.Option.Body)
	if err != nil {
		return err
	}

	publisher.NotifyPublish(func(r amqp.Confirmation) {
		if !r.Confirmation.Ack {
			Logrus(cons.ERROR, "Failed message delivery to: %s", req.Option.QueueName)
			return
		}

		Logrus(cons.INFO, "Success message delivery to: %s", req.Option.QueueName)
		return
	})

	err = publisher.PublishWithContext(p.ctx, bodyByte, []string{req.Option.QueueName},
		amqp.WithPublishOptionsPersistentDelivery,
		amqp.WithPublishOptionsExchange(req.Option.ExchangeName),
		amqp.WithPublishOptionsContentType(req.Option.ContentType),
		amqp.WithPublishOptionsTimestamp(req.Option.Timestamp),
		amqp.WithPublishOptionsHeaders(req.Option.Args),
	)

	if err != nil {
		return err
	}

	return nil
}

func (p rabbitmq) Consumer(req dto.Request[dto.RabbitOptions], callback func(d amqp.Delivery) (action amqp.Action)) {
	if req.Option.ConsumerID == "" {
		req.Option.ConsumerID = shortuuid.New()
	}

	if req.Option.Concurrency < 1 {
		req.Option.Concurrency = runtime.NumCPU() / 2
	}

	if req.Option.Prefetch < 1 {
		req.Option.Concurrency = 5
	}

	consumer, err := amqp.NewConsumer(p.rabbitmq, callback, req.Option.QueueName,
		amqp.WithConsumerOptionsExchangeName(req.Option.ExchangeName),
		amqp.WithConsumerOptionsExchangeKind(req.Option.ExchangeType),
		amqp.WithConsumerOptionsBinding(amqp.Binding{
			RoutingKey: req.Option.QueueName,
			BindingOptions: amqp.BindingOptions{
				Declare: true,
				NoWait:  true,
				Args:    req.Option.Args,
			},
		}),
		amqp.WithConsumerOptionsExchangeDurable,
		amqp.WithConsumerOptionsQueueDurable,
		amqp.WithConsumerOptionsConsumerName(req.Option.ConsumerID),
		amqp.WithConsumerOptionsConsumerAutoAck(req.Option.Ack),
		amqp.WithConsumerOptionsConcurrency(req.Option.Concurrency),
		amqp.WithConsumerOptionsQOSPrefetch(req.Option.Prefetch),
		amqp.WithConsumerOptionsQueueArgs(req.Option.Args),
		amqp.WithConsumerOptionsLogging,
	)
	defer p.closeConnection(nil, consumer, p.rabbitmq)

	if err != nil {
		Logrus(cons.ERROR, err)
		return
	}

	return
}

func (h *rabbitmq) closeConnection(publisher *amqp.Publisher, consumer *amqp.Consumer, connection *amqp.Conn) {
	defer h.recovery()

	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM, syscall.SIGABRT, syscall.SIGUSR1)

	go func() {
		<-closeChan

		if consumer != nil {
			consumer.Close()
		}

		if publisher != nil {
			publisher.Close()
		}

		if connection != nil {
			connection.Close()
		}
	}()
}

func (p rabbitmq) recovery() {
	if err := recover(); err != nil {
		return
	}
}
