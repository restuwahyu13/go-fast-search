package service

import (
	"context"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
)

type usersService struct {
	env dto.Request[dto.Environtment]
	db  *bun.DB
	rds *redis.Client
}

func NewUsersService(options dto.ServiceOptions) inf.IUsersService {
	return usersService{env: options.ENV, db: options.DB, rds: options.RDS}
}

func (s usersService) Ping(ctx context.Context) (res opt.Response) {
	res.StatCode = http.StatusOK
	res.Message = "Ping!"

	return
}

func (s usersService) CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) (res opt.Response) {
	res.StatCode = http.StatusOK
	res.Message = "Ping!"

	return
}
