package dto

import opt "github.com/restuwahyu13/go-fast-search/shared/output"

type Config struct {
	ENV              string `env:"GO_ENV" mapstructure:"GO_ENV"`
	PORT             string `env:"PORT" mapstructure:"PORT"`
	INBOUND_SIZE     int    `env:"INBOUND_SIZE" mapstructure:"INBOUND_SIZE"`
	PG_DSN           string `env:"PG_DSN" mapstructure:"PG_DSN"`
	REDIS_CSN        string `env:"REDIS_CSN" mapstructure:"REDIS_CSN"`
	JWT_SECRET_KEY   string `env:"JWT_SECRET_KEY" mapstructure:"JWT_SECRET_KEY"`
	JWT_EXPIRED      int    `env:"JWT_EXPIRED" mapstructure:"JWT_EXPIRED"`
	RABBITMQ_QSN     string `env:"RABBITMQ_QSN" mapstructure:"RABBITMQ_QSN"`
	RABBITMQ_VSN     string `env:"RABBITMQ_VSN" mapstructure:"RABBITMQ_VSN"`
	MEILI_DSN        string `env:"MEILI_DSN" mapstructure:"MEILI_DSN"`
	MEILI_MASTER_KEY string `env:"MEILI_MASTER_KEY" mapstructure:"MEILI_MASTER_KEY"`
}

type (
	Environtment struct {
		APP         opt.Application
		REDIS       opt.Redis
		POSTGRES    opt.Postgres
		JWT         opt.Jwt
		RABBITMQ    opt.RabbitMQ
		MEILISEARCH opt.MeiliSearch
	}
)
