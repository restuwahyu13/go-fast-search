package dto

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/wagslane/go-rabbitmq"
)

type (
	ServiceOptions struct {
		ENV  Request[Environtment]
		DB   *bun.DB
		RDS  *redis.Client
		AMQP *rabbitmq.Conn
		MLS  meilisearch.ServiceManager
	}

	UsecaseOptions[T any] struct {
		SERVICE T
	}

	ControllerOptions[T any] struct {
		USECASE T
	}

	RouteOptions[T any] struct {
		ENV        Request[Environtment]
		RDS        *redis.Client
		ROUTER     chi.Router
		CONTROLLER T
	}

	WorkerOptions struct {
		CTX  context.Context
		ENV  Request[Environtment]
		DB   *bun.DB
		RDS  *redis.Client
		AMQP *rabbitmq.Conn
		MLS  meilisearch.ServiceManager
	}

	ModuleOptions struct {
		ENV    Request[Environtment]
		DB     *bun.DB
		RDS    *redis.Client
		AMQP   *rabbitmq.Conn
		MLS    meilisearch.ServiceManager
		ROUTER chi.Router
	}
)
