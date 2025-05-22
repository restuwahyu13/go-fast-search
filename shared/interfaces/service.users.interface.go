package inf

import (
	"context"
	"net/http"

	"github.com/uptrace/bun"

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
		Ping(ctx context.Context) opt.Response
	}

	IUsersException interface{}

	IUsersUsecase interface {
		Ping(ctx context.Context) opt.Response
	}

	IUsersController interface {
		Ping(rw http.ResponseWriter, r *http.Request)
	}
)
