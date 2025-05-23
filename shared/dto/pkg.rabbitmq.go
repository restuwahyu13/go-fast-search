package dto

import (
	"time"

	"github.com/wagslane/go-rabbitmq"
)

type (
	RabbitOptions struct {
		ExchangeName string
		ExchangeType string
		QueueName    string
		Ack          bool
		Concurrency  int
		ConsumerID   string
		Args         rabbitmq.Table
		Body         any
		ContentType  string
		Timestamp    time.Time
		Prefetch     int
	}

	RabbitDeadLetterQueueOptions struct {
		Exchange     string
		ExchangeType string
		Queue        string
		Body         MeiliSearchDocuments[map[string]any]
		Secret       any
		Unknown      bool
		Error        error
	}
)
