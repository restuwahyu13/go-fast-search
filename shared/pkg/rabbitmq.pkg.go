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
	env      dto.Environtment
	rabbitmq *amqp.Conn
}

func NewRabbitMQ(ctx context.Context, env dto.Environtment, con *amqp.Conn) inf.IRabbitMQ {
	return rabbitmq{ctx: ctx, env: env, rabbitmq: con}
}

func (p rabbitmq) Publisher(req dto.Request[dto.RabbitOptions]) error {
	parser := helper.NewParser()

	if req.Option.ContentType != "" {
		req.Option.ContentType = "application/json"
	}

	if req.Option.Timestamp.Sub(time.Now()).Seconds() < 0 {
		req.Option.Timestamp = time.Now().Local()
	}

	publisher, err := amqp.NewPublisher(p.rabbitmq,
		amqp.WithPublisherOptionsExchangeName(req.Option.ExchangeName),
		amqp.WithPublisherOptionsExchangeKind(req.Option.ExchangeType),
		amqp.WithPublisherOptionsExchangeDeclare,
		amqp.WithPublisherOptionsExchangeDurable,
		amqp.WithPublisherOptionsExchangeNoWait,
		amqp.WithPublisherOptionsLogging,
	)

	defer p.closeConnection(publisher, nil, p.rabbitmq)
	if err != nil {
		return err
	}

	bodyByte, err := parser.Marshal(req.Option.Body)
	if err != nil {
		return err
	}

	err = publisher.Publish(bodyByte, []string{req.Option.QueueName},
		amqp.WithPublishOptionsPersistentDelivery,
		amqp.WithPublishOptionsExchange(req.Option.ExchangeName),
		amqp.WithPublishOptionsContentType(req.Option.ContentType),
		amqp.WithPublishOptionsTimestamp(req.Option.Timestamp),
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

	consumer, err := amqp.NewConsumer(p.rabbitmq, callback, req.Option.QueueName,
		amqp.WithConsumerOptionsExchangeName(req.Option.ExchangeName),
		amqp.WithConsumerOptionsExchangeKind(req.Option.ExchangeType),
		amqp.WithConsumerOptionsBinding(amqp.Binding{
			RoutingKey: req.Option.QueueName,
			BindingOptions: amqp.BindingOptions{
				Declare: true,
				NoWait:  false,
				Args:    req.Option.Args,
			},
		}),
		amqp.WithConsumerOptionsExchangeDurable,
		amqp.WithConsumerOptionsQueueDurable,
		amqp.WithConsumerOptionsConsumerName(req.Option.ConsumerID),
		amqp.WithConsumerOptionsConsumerAutoAck(req.Option.Ack),
		amqp.WithConsumerOptionsConcurrency(req.Option.Concurrency),
		amqp.WithConsumerOptionsLogging,
	)

	if err != nil {
		Logrus(cons.ERROR, err)
		p.closeConnection(nil, consumer, p.rabbitmq)
		return
	}
}

func (p rabbitmq) closeConnection(publisher *amqp.Publisher, consumer *amqp.Consumer, connection *amqp.Conn) {
	p.recovery()

	if publisher != nil && consumer != nil && connection != nil {
		publisher.Close()
		consumer.Close()
	} else if publisher != nil && consumer == nil && connection != nil {
		publisher.Close()
	} else if publisher == nil && consumer != nil && connection != nil {
		consumer.Close()
	} else {
		closeChan := make(chan os.Signal, 1)
		signal.Notify(closeChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM, syscall.SIGABRT, syscall.SIGUSR1)

		for {
			select {
			case <-closeChan:
				publisher.Close()
				consumer.Close()
				connection.Close()
			default:
				return
			}
		}
	}
}

func (p rabbitmq) recovery() {
	if err := recover(); err != nil {
		return
	}
}
