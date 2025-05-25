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
		Find() *bun.SelectQuery
		FindOne() *bun.SelectQuery
		Insert(entitie entitie.UsersEntitie, column string, dest ...any) error
		Update(entitie entitie.UsersEntitie, column string, dest ...any) error
		Delete(id string, dest any) error
	}

	IUsersMeiliSearchRepositorie interface {
		Search(query string, filter *meilisearch.SearchRequest) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error)
		Find(filter *meilisearch.DocumentsQuery) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error)
		FindOne(id string, filter *meilisearch.DocumentQuery) (*entitie.UsersDocument, error)
		Insert(value any) error
		Update(id string, value any) error
		Delete(id string) error
		BulkInsert(value any) error
		BulkUpdate(value any) error
		UpdateFilterableAttributes(attributes ...string) error
		UpdateSearchableAttributes(attributes ...string) error
		UpdateSortableAttributes(attributes ...string) error
		UpdateDisplayedAttributes(attributes ...string) error
		ListUsersDocuments(req dto.Request[dto.MeiliSearchDocumentsQuery]) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error)
	}

	IUsersService interface {
		CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) (res opt.Response)
		UpdateUsers(ctx context.Context, req dto.Request[dto.UpdateUsersDTO]) (res opt.Response)
		FindAllUsers(ctx context.Context, req dto.Request[dto.MeiliSearchDocumentsQuery]) (res opt.Response)
	}

	IUsersException interface {
		CreateUsers(key string) string
		UpdateUsers(key string) string
	}

	IUsersUsecase interface {
		CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) opt.Response
		UpdateUsers(ctx context.Context, req dto.Request[dto.UpdateUsersDTO]) opt.Response
		FindAllUsers(ctx context.Context, req dto.Request[dto.MeiliSearchDocumentsQuery]) opt.Response
	}

	IUsersController interface {
		CreateUsers(rw http.ResponseWriter, r *http.Request)
		UpdateUsers(rw http.ResponseWriter, r *http.Request)
		FindAllUsers(rw http.ResponseWriter, r *http.Request)
	}
)
