package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"

	entitie "github.com/restuwahyu13/go-fast-search/domain/entities"
	repo "github.com/restuwahyu13/go-fast-search/domain/repositories"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
)

type usersService struct {
	env dto.Request[dto.Environtment]
	db  *bun.DB
	rds *redis.Client
	mls meilisearch.ServiceManager
}

func NewUsersService(options dto.ServiceOptions) inf.IUsersService {
	return usersService{env: options.ENV, db: options.DB, rds: options.RDS, mls: options.MLS}
}

func (s usersService) Ping(ctx context.Context) (res opt.Response) {
	res.StatCode = http.StatusOK
	res.Message = "Ping!"

	return
}

func (s usersService) CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) (res opt.Response) {
	usersRepositorie := repo.NewUsersRepositorie(ctx, s.db)
	usersEntitie := entitie.UsersEntitie{}

	err := usersRepositorie.FindOne(&usersEntitie).Column("id").
		Where("deleted_at IS NULL").
		Where("name = ?", req.Body.Name).
		Where("email = ?", req.Body.Email).
		Where("phone = ?", req.Body.Phone).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		res.StatCode = http.StatusInternalServerError
		res.ErrMsg = err.Error()

		return

	} else if err != sql.ErrNoRows {
		res.StatCode = http.StatusConflict
		res.ErrMsg = "User already exists in our system"

		return

	}

	usersEntitie.Name = req.Body.Name
	usersEntitie.Email = req.Body.Email
	usersEntitie.Phone = req.Body.Phone
	usersEntitie.DateOfBirth = req.Body.DateOfBirth
	usersEntitie.Address = req.Body.Address
	usersEntitie.City = req.Body.City
	usersEntitie.State = req.Body.State
	usersEntitie.Direction = req.Body.Direction
	usersEntitie.Country = req.Body.Country
	usersEntitie.PostalCode = req.Body.PostalCode

	if err := usersRepositorie.Insert(&usersEntitie, nil); err != nil {
		if err != cons.NO_ROWS_AFFECTED {
			res.StatCode = http.StatusInternalServerError
			res.ErrMsg = err.Error()

			return
		}

		res.StatCode = http.StatusPreconditionFailed
		res.ErrMsg = "Failed to create new users"

		return
	}

	res.StatCode = http.StatusOK
	res.Message = "Success to create new users"

	return
}

func (s usersService) UpdateUsers(ctx context.Context, req dto.Request[dto.UpdateUsersDTO]) (res opt.Response) {
	usersRepositorie := repo.NewUsersRepositorie(ctx, s.db)
	usersEntitie := entitie.UsersEntitie{}

	err := usersRepositorie.FindOne(&usersEntitie).Column("id").
		Where("deleted_at IS NULL").
		Where("id = ?", req.Param.ID).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		res.StatCode = http.StatusInternalServerError
		res.ErrMsg = err.Error()

		return

	} else if err == sql.ErrNoRows {
		res.StatCode = http.StatusConflict
		res.ErrMsg = "User is not exists in our system"

		return

	}

	usersEntitie.Name = req.Body.Name
	usersEntitie.Email = req.Body.Email
	usersEntitie.Phone = req.Body.Phone
	usersEntitie.DateOfBirth = req.Body.DateOfBirth
	usersEntitie.Address = req.Body.Address
	usersEntitie.City = req.Body.City
	usersEntitie.State = req.Body.State
	usersEntitie.Direction = req.Body.Direction
	usersEntitie.Country = req.Body.Country
	usersEntitie.PostalCode = req.Body.PostalCode

	if err := usersRepositorie.Update(&usersEntitie, nil); err != nil {
		if err != cons.NO_ROWS_AFFECTED {
			res.StatCode = http.StatusInternalServerError
			res.ErrMsg = err.Error()

			return
		}

		res.StatCode = http.StatusPreconditionFailed
		res.ErrMsg = "Failed to update new users"

		return
	}

	res.StatCode = http.StatusOK
	res.Message = "Success to update users"
	res.Data = usersEntitie

	return
}

func (s usersService) FindAllUsers(ctx context.Context) (res opt.Response) {
	usersRepositorie := repo.NewUsersMeilisearchRepositorie(ctx, s.mls)

	usersDocResult, err := usersRepositorie.Find("users", &meilisearch.DocumentsQuery{})
	if err != nil {
		res.StatCode = http.StatusInternalServerError
		res.ErrMsg = err.Error()

		return
	}

	res.StatCode = http.StatusOK
	res.Message = "Success"
	res.Data = usersDocResult

	return
}
