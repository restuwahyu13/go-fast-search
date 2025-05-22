package config

import (
	"os"

	genv "github.com/caarlos0/env"
	"github.com/spf13/viper"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
)

func NewEnvirontment(name, path, ext string) (*opt.Environtment, error) {
	cfg := dto.Config{}

	if _, ok := os.LookupEnv("GO_ENV"); !ok {
		viper.SetConfigName(name)
		viper.SetConfigType(ext)
		viper.AddConfigPath(path)
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}

		if err := viper.Unmarshal(&cfg); err != nil {
			return nil, err
		}
	} else {
		if err := genv.Parse(&cfg); err != nil {
			return nil, err
		}
	}

	return &opt.Environtment{
		APP: opt.Application{
			ENV:          cfg.ENV,
			PORT:         cfg.PORT,
			INBOUND_SIZE: cfg.INBOUND_SIZE,
		},
		REDIS: opt.Redis{
			URL: cfg.REDIS_CSN,
		},
		POSTGRES: opt.Postgres{
			URL: cfg.PG_DSN,
		},
		JWT: opt.Jwt{
			SECRET:  cfg.JWT_SECRET_KEY,
			EXPIRED: cfg.JWT_EXPIRED,
		},
		RABBITMQ: opt.RabbitMQ{
			URL: cfg.RABBITMQ_QSN,
			VSN: cfg.RABBITMQ_VSN,
		},
		MEILISEARCH: opt.MeiliSearch{
			URL: cfg.MEILI_DSN,
			KEY: cfg.MEILI_MASTER_KEY,
		},
	}, nil
}
