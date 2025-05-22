package inf

import (
	"github.com/wagslane/go-rabbitmq"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
)

type IRabbitMQ interface {
	Publisher(req dto.Request[dto.RabbitOptions]) error
	Consumer(req dto.Request[dto.RabbitOptions], callback func(d rabbitmq.Delivery) (action rabbitmq.Action))
}
