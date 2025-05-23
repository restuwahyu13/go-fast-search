package inf

import (
	"context"
	"net/http"

	"github.com/uptrace/bun"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
)

type (
	IUsersRepositorie interface {
		Find(dest any) *bun.SelectQuery
		FindOne(dest any) *bun.SelectQuery
		Create(dest any) *bun.InsertQuery
		Update(dest any) *bun.UpdateQuery
		Delete(dest any) *bun.DeleteQuery
	}

	IUsersService interface {
		Ping(ctx context.Context) (res opt.Response)
		CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) (res opt.Response)
	}

	IUsersException interface{}

	IUsersUsecase interface {
		Ping(ctx context.Context) opt.Response
		CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) opt.Response
	}

	IUsersController interface {
		Ping(rw http.ResponseWriter, r *http.Request)
		CreateUsers(rw http.ResponseWriter, r *http.Request)
	}
)
