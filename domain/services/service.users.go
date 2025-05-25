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

	err := usersRepositorie.FindOne().Column("id").
		Where("deleted_at IS NULL").
		Where("name = ?", req.Body.Name).
		Where("email = ?", req.Body.Email).
		Where("phone = ?", req.Body.Phone).
		Scan(ctx, &usersEntitie)

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

	if err := usersRepositorie.Insert(usersEntitie, "id, created_at", &usersEntitie.ID, &usersEntitie.CreatedAt); err != nil {
		if err != cons.NO_ROWS_AFFECTED {
			res.StatCode = http.StatusInternalServerError
			res.ErrMsg = err.Error()

			return
		}

		res.StatCode = http.StatusPreconditionFailed
		res.ErrMsg = "Failed to create new users"

		return
	}

	createdAtUnix, err := helper.TimeStampToUnix(usersEntitie.CreatedAt.Format(time.RFC3339))
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	updatedAtUnix, err := helper.TimeStampToUnix(usersEntitie.UpdatedAt.Time.Format(time.RFC3339))
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	deletedAtUnix, err := helper.TimeStampToUnix(usersEntitie.DeletedAt.Time.Format(time.RFC3339))
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	usersDocEntitie := entitie.UsersDocument{}
	usersDocEntitie.ID = usersEntitie.ID
	usersDocEntitie.Name = usersEntitie.Name
	usersDocEntitie.Email = usersEntitie.Email
	usersDocEntitie.Phone = usersEntitie.Phone
	usersDocEntitie.DateOfBirth = usersEntitie.DateOfBirth
	usersDocEntitie.Age = usersEntitie.Age
	usersDocEntitie.Address = usersEntitie.Address
	usersDocEntitie.City = usersEntitie.City
	usersDocEntitie.State = usersEntitie.State
	usersDocEntitie.Direction = usersEntitie.Direction
	usersDocEntitie.Country = usersEntitie.Country
	usersDocEntitie.PostalCode = usersEntitie.PostalCode
	usersDocEntitie.CreatedAt = createdAtUnix
	usersDocEntitie.UpdatedAt = updatedAtUnix
	usersDocEntitie.DeletedAt = deletedAtUnix

	amqp := pkg.NewRabbitMQ(ctx, s.amqp)

	if err := helper.MeiliSearchPublisher[entitie.UsersDocument](amqp, s.env.Config.RABBITMQ.SECRET, nil, usersDocEntitie, cons.FALSE, cons.INSERT); err != nil {
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

	err := usersRepositorie.FindOne().Column("id").
		Where("deleted_at IS NULL").
		Where("id = ?", req.Body.ID).
		Scan(ctx, &usersEntitie)

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

	if err := usersRepositorie.Update(usersEntitie, "*", &usersEntitie); err != nil {
		if err != cons.NO_ROWS_AFFECTED {
			res.StatCode = http.StatusInternalServerError
			res.ErrMsg = err.Error()

			return
		}

		res.StatCode = http.StatusPreconditionFailed
		res.ErrMsg = "Failed to update new users"

		return
	}

	createdAtUnix, err := helper.TimeStampToUnix(usersEntitie.CreatedAt.Format(time.RFC3339))
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	updatedAtUnix, err := helper.TimeStampToUnix(usersEntitie.UpdatedAt.Time.Format(time.RFC3339))
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	deletedAtUnix, err := helper.TimeStampToUnix(usersEntitie.DeletedAt.Time.Format(time.RFC3339))
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	usersDocEntitie := entitie.UsersDocument{}
	usersDocEntitie.ID = usersEntitie.ID
	usersDocEntitie.Name = usersEntitie.Name
	usersDocEntitie.Email = usersEntitie.Email
	usersDocEntitie.Phone = usersEntitie.Phone
	usersDocEntitie.DateOfBirth = usersEntitie.DateOfBirth
	usersDocEntitie.Age = usersEntitie.Age
	usersDocEntitie.Address = usersEntitie.Address
	usersDocEntitie.City = usersEntitie.City
	usersDocEntitie.State = usersEntitie.State
	usersDocEntitie.Direction = usersEntitie.Direction
	usersDocEntitie.Country = usersEntitie.Country
	usersDocEntitie.PostalCode = usersEntitie.PostalCode
	usersDocEntitie.CreatedAt = createdAtUnix
	usersDocEntitie.UpdatedAt = updatedAtUnix
	usersDocEntitie.DeletedAt = deletedAtUnix

	amqp := pkg.NewRabbitMQ(ctx, s.amqp)

	if err := helper.MeiliSearchPublisher[entitie.UsersDocument](amqp, s.env.Config.RABBITMQ.SECRET, req.Body.ID, usersDocEntitie, cons.FALSE, cons.UPDATE); err != nil {
		res.StatCode = http.StatusInternalServerError
		res.ErrMsg = err.Error()

		return
	}

	res.StatCode = http.StatusOK
	res.Message = "Success to update users"
	res.Data = usersEntitie

	return
}

func (s usersService) FindAllUsers(ctx context.Context, req dto.Request[dto.MeiliSearchDocumentsQuery]) (res opt.Response) {
	if req.Query.Limit < 1 {
		req.Query.Limit = 10
	}

	if req.Query.Page < 1 {
		req.Query.Page = 1
	}

	if req.Query.SearchBy == "" {
		req.Query.SearchBy = "name,email,phone"
	}

	if req.Query.SortBy == "" {
		req.Query.SortBy = "created_at"
	}

	if req.Query.Sort == "" {
		req.Query.Sort = "desc"
	}

	req.Query.Page = (req.Query.Page - 1) * req.Query.Limit

	usersRepositorie := repo.NewUsersMeilisearchRepositorie(ctx, s.mls)
	resultUsersDocuments, err := usersRepositorie.ListUsersDocuments(req)
	if err != nil {
		res.StatCode = http.StatusInternalServerError
		res.ErrMsg = err.Error()

		return
	}

	res.StatCode = http.StatusOK
	res.Message = "Success"
	res.Data = resultUsersDocuments

	return
}
