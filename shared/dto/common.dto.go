package dto

import (
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

type (
	ServiceOptions struct {
		ENV Request[Environtment]
		DB  *bun.DB
		RDS *redis.Client
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

	ModuleOptions struct {
		ENV    Request[Environtment]
		DB     *bun.DB
		RDS    *redis.Client
		ROUTER chi.Router
	}
)
