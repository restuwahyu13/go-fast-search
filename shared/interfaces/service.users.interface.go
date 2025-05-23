package inf

import (
	"context"
	"net/http"

	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"

	entitie "github.com/restuwahyu13/go-fast-search/domain/entities"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
)

type (
	IUsersRepositorie interface {
		Find(dest any) *bun.SelectQuery
		FindOne(dest any) *bun.SelectQuery
		Insert(entitie any, dest any) error
		Update(entitie any, dest any) error
		Delete(id string, dest any) error
	}

	IUsersMeiliSearchRepositorie interface {
		Find(filter *meilisearch.DocumentsQuery) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error)
		FindOne(id string, filter *meilisearch.DocumentQuery) (*entitie.UsersDocument, error)
		Insert(value any) error
		Update(id string, value any) error
		Delete(id string) error
		BulkInsert(value any) error
		BulkUpdate(value any) error
	}

	IUsersService interface {
		Ping(ctx context.Context) (res opt.Response)
		CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) (res opt.Response)
		UpdateUsers(ctx context.Context, req dto.Request[dto.UpdateUsersDTO]) (res opt.Response)
		FindAllUsers(ctx context.Context) (res opt.Response)
	}

	IUsersException interface{}

	IUsersUsecase interface {
		Ping(ctx context.Context) opt.Response
		CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) opt.Response
		UpdateUsers(ctx context.Context, req dto.Request[dto.UpdateUsersDTO]) opt.Response
		FindAllUsers(ctx context.Context) opt.Response
	}

	IUsersController interface {
		Ping(rw http.ResponseWriter, r *http.Request)
		CreateUsers(rw http.ResponseWriter, r *http.Request)
		UpdateUsers(rw http.ResponseWriter, r *http.Request)
		FindAllUsers(rw http.ResponseWriter, r *http.Request)
	}
)
