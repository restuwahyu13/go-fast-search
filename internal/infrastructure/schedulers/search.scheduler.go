package scheduler

import (
	"context"
	"fmt"
	"time"

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
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type searchScheduler struct {
	ctx  context.Context
	env  dto.Request[dto.Environtment]
	db   *bun.DB
	rds  *redis.Client
	amqp *rabbitmq.Conn
	mls  meilisearch.ServiceManager
}

func NewSearchScheduler(options dto.SchedulerOptions) inf.ISearchScheduler {
	return searchScheduler{ctx: options.CTX, env: options.ENV, db: options.DB, rds: options.RDS, amqp: options.AMQP, mls: options.MLS}
}

func (s searchScheduler) searchHandler(rds inf.IRedis) {
	key := "WORKER:SEARCH:CDC"

	isExists, err := rds.Exists(key)
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	if isExists > 0 {
		result, err := rds.Get(key)
		if err != nil {
			pkg.Logrus(cons.ERROR, err)
			return
		}

		start_at := string(result)
		limit := 100

		usersRepositorie := repo.NewUsersRepositorie(s.ctx, s.db)
		usersEntities := []entitie.UsersEntitie{}

		err = usersRepositorie.Find().Column("*").
			WhereGroup("AND", func(sqlb *bun.SelectQuery) *bun.SelectQuery {
				sqlb.Where("deleted_at IS NULL AND updated_at IS NULL AND created_at > ?", start_at)
				sqlb.WhereOr("deleted_at IS NULL AND updated_at > ?", start_at)
				return sqlb
			}).
			Order("created_at DESC").
			Limit(limit).Scan(s.ctx, &usersEntities)

		if err != nil {
			pkg.Logrus(cons.ERROR, err)
			return
		}

		pkg.Logrus(cons.INFO, "Total data %d sync to meilisearch", len(usersEntities))

		if len(usersEntities) > 0 {
			usersRepositorie := repo.NewUsersMeilisearchRepositorie(s.ctx, s.mls)

			for _, userEntity := range usersEntities {
				createdAtUnix, err := helper.TimeStampToUnix(userEntity.CreatedAt.Format(time.RFC3339))
				if err != nil {
					pkg.Logrus(cons.ERROR, err)
					return
				}

				updatedAtUnix, err := helper.TimeStampToUnix(userEntity.UpdatedAt.Time.Format(time.RFC3339))
				if err != nil {
					pkg.Logrus(cons.ERROR, err)
					return
				}

				deletedAtUnix, err := helper.TimeStampToUnix(userEntity.DeletedAt.Time.Format(time.RFC3339))
				if err != nil {
					pkg.Logrus(cons.ERROR, err)
					return
				}

				usersDocEntitie := entitie.UsersDocument{}
				usersDocEntitie.ID = userEntity.ID
				usersDocEntitie.Name = userEntity.Name
				usersDocEntitie.Email = userEntity.Email
				usersDocEntitie.Phone = userEntity.Phone
				usersDocEntitie.DateOfBirth = userEntity.DateOfBirth
				usersDocEntitie.Age = userEntity.Age
				usersDocEntitie.Address = userEntity.Address
				usersDocEntitie.City = userEntity.City
				usersDocEntitie.State = userEntity.State
				usersDocEntitie.Direction = userEntity.Direction
				usersDocEntitie.Country = userEntity.Country
				usersDocEntitie.PostalCode = userEntity.PostalCode
				usersDocEntitie.CreatedAt = createdAtUnix
				usersDocEntitie.UpdatedAt = updatedAtUnix
				usersDocEntitie.DeletedAt = deletedAtUnix

				defaultTimeUnix, err := helper.TimeStampToUnix(time.Time{}.Format(time.RFC3339))
				if err != nil {
					pkg.Logrus(cons.ERROR, err)
					return
				}

				cdcTimeUnix, err := helper.TimeStampToUnix(start_at)
				if err != nil {
					pkg.Logrus(cons.ERROR, err)
					return
				}

				createdAtFilter := fmt.Sprintf("deleted_at = %d AND updated_at = %d AND created_at > %d", defaultTimeUnix, defaultTimeUnix, cdcTimeUnix)
				updatedAtFilter := fmt.Sprintf("deleted_at = %d AND updated_at > %d", defaultTimeUnix, cdcTimeUnix)

				filter := fmt.Sprintf("(%s) OR (%s)", createdAtFilter, updatedAtFilter)
				attributes := []string{"deleted_at", "created_at", "updated_at"}

				filterSearch := meilisearch.SearchRequest{Filter: filter, Limit: 1}
				usersDoc, err := usersRepositorie.Search("", attributes, &filterSearch)

				if err != nil {
					pkg.Logrus(cons.ERROR, err)
					return
				}

				if usersDoc.Results != nil {
					if err := usersRepositorie.Update(usersDocEntitie.ID, usersDocEntitie); err != nil {
						pkg.Logrus(cons.ERROR, err)
						return
					}
				} else {
					if err := usersRepositorie.Insert(usersDocEntitie); err != nil {
						pkg.Logrus(cons.ERROR, err)
						return
					}
				}
			}

			pkg.Logrus(cons.INFO, "Total data %d sync to meilisearch success", len(usersEntities))
		}
	}
}

func (s searchScheduler) Run() {
	cron := pkg.NewCron()

	crontime := cons.Every15Seconds
	now := time.Now().Format(cons.DATE_TIME_FORMAT)

	rds, err := pkg.NewRedis(s.ctx, s.rds)
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	key := "SCHEDULER:SEARCH:BREAK"
	value := 1
	sync := 10

	sch, _, err := cron.Handler("search scheduler", crontime, func() {
		pkg.Logrus(cons.INFO, fmt.Sprintf("Search scheduler is running %s - and execute at %s", now, crontime))

		result, err := rds.IncrBy(key, value)
		if err != nil {
			pkg.Logrus(cons.ERROR, err)
			return
		}

		if result >= sync {
			breakTime := time.Duration(time.Second * 180)
			pkg.Logrus(cons.INFO, fmt.Sprintf("Search scheduler break when equal max %d, running again after %d minute is over", sync, int64(breakTime.Minutes())))

			ttl, err := rds.TTL(key, int(breakTime.Seconds()))
			if err != nil {
				pkg.Logrus(cons.ERROR, err)
				return
			}

			if ttl < 1 {
				if err := rds.SetEx(key, breakTime, result); err != nil {
					pkg.Logrus(cons.ERROR, err)
					return
				}

			} else if ttl > 1 && ttl < 3 {
				if err := rds.Set(key, value); err != nil {
					pkg.Logrus(cons.ERROR, err)
					return
				}

			}
		}

		if result <= sync {
			s.searchHandler(rds)
		}
	})

	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	sch.Start()
}
