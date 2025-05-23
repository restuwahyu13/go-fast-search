package cons

const (
	X_RABBIT_SECRET        = "x-rabbit-secret"
	X_RABBIT_UNKNOWN       = "x-rabbit-unknown"
	X_RABBIT_QUEUE         = "x-rabbit-queue"
	X_RABBIT_EXCHANGE      = "x-rabbit-exchange"
	X_RABBIT_EXCHANGE_TYPE = "x-rabbit-exchange-type"
	X_MESSAGE_TTL          = "x-message-ttl"
)

const (
	EXCHANGE_TYPE_DIRECT = "direct"
	EXCHANGE_TYPE_FANOUT = "fanout"
	EXCHANGE_TYPE_TOPIC  = "topic"
	EXCHANGE_TYPE_HEADER = "header"
)

const (
	EXCHANGE_NAME_DIRECT            = "amq.direct"
	EXCHANGE_NAME_TOPIC             = "amq.topic"
	EXCHANGE_NAME_SEARCH            = "amqp.worker"
	EXCHANGE_NAME_DEAD_LETTER_QUEUE = "amqp.worker.dlq"
)

const (
	QUEUE_NAME_SEARCH            = "worker.search"
	QUEUE_NAME_DEAD_LETTER_QUEUE = "worker.dlq"
)
