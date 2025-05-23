package service

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/wagslane/go-rabbitmq"

	entitie "github.com/restuwahyu13/go-fast-search/domain/entities"
	repo "github.com/restuwahyu13/go-fast-search/domain/repositories"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type usersService struct {
	env  dto.Request[dto.Environtment]
	db   *bun.DB
	rds  *redis.Client
	amqp *rabbitmq.Conn
	mls  meilisearch.ServiceManager
}

func NewUsersService(options dto.ServiceOptions) inf.IUsersService {
	return usersService{env: options.ENV, db: options.DB, rds: options.RDS, amqp: options.AMQP, mls: options.MLS}
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
	usersEntitie.Age = req.Body.Age
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

	amqp := pkg.NewRabbitMQ(ctx, s.amqp)

	if err := helper.MeiliSearchPublisher[entitie.UsersEntitie](amqp, s.env.Config.RABBITMQ.SECRET, nil, usersEntitie, cons.FALSE, cons.INSERT); err != nil {
		res.StatCode = http.StatusInternalServerError
		res.ErrMsg = err.Error()

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
		Where("id = ?", req.Body.ID).
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
	usersEntitie.Age = req.Body.Age
	usersEntitie.Address = req.Body.Address
	usersEntitie.City = req.Body.City
	usersEntitie.State = req.Body.State
	usersEntitie.Direction = req.Body.Direction
	usersEntitie.Country = req.Body.Country
	usersEntitie.PostalCode = req.Body.PostalCode
	usersEntitie.UpdatedAt = zero.TimeFrom(time.Now())

	if err := usersRepositorie.Update(&usersEntitie, &usersEntitie); err != nil {
		if err != cons.NO_ROWS_AFFECTED {
			res.StatCode = http.StatusInternalServerError
			res.ErrMsg = err.Error()

			return
		}

		res.StatCode = http.StatusPreconditionFailed
		res.ErrMsg = "Failed to update new users"

		return
	}

	amqp := pkg.NewRabbitMQ(ctx, s.amqp)

	if err := helper.MeiliSearchPublisher[entitie.UsersEntitie](amqp, s.env.Config.RABBITMQ.SECRET, req.Body.ID, usersEntitie, cons.FALSE, cons.UPDATE); err != nil {
		res.StatCode = http.StatusInternalServerError
		res.ErrMsg = err.Error()

		return
	}

	res.StatCode = http.StatusOK
	res.Message = "Success to update users"
	res.Data = usersEntitie

	return
}

func (s usersService) FindAllUsers(ctx context.Context) (res opt.Response) {
	usersRepositorie := repo.NewUsersMeilisearchRepositorie(ctx, s.mls)

	usersDocResult, err := usersRepositorie.Find(&meilisearch.DocumentsQuery{})
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
