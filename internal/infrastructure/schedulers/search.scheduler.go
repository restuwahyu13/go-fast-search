package scheduler

import (
	"context"
	"fmt"
	"sync"
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

func (s searchScheduler) findAllUsers(wg *sync.WaitGroup, startAt string, usersEntitiesChan chan []entitie.UsersEntitie, errChan chan error) (*sync.WaitGroup, string, <-chan []entitie.UsersEntitie, chan error) {
	defer wg.Done()

	usersRepositorie := repo.NewUsersRepositorie(s.ctx, s.db)
	usersEntities := []entitie.UsersEntitie{}

	limit := 1000

	err := usersRepositorie.Find().Column("*").
		Where("deleted_at IS NULL").
		WhereGroup(cons.AND, func(sqlb *bun.SelectQuery) *bun.SelectQuery {
			sqlb.Where("updated_at IS NULL AND created_at > ?", startAt)
			sqlb.WhereOr("updated_at > ?", startAt)

			return sqlb
		}).
		Order("created_at DESC").
		Limit(limit).Scan(s.ctx, &usersEntities)

	if err != nil {
		errChan <- err
		return wg, "", usersEntitiesChan, errChan
	}

	usersEntitiesChan <- usersEntities

	return wg, startAt, usersEntitiesChan, errChan
}

func (s searchScheduler) updateUsers(wg *sync.WaitGroup, start_at string, usersEntitiesChan <-chan []entitie.UsersEntitie, errChan chan error) {
	defer wg.Done()

	usersEntities := <-usersEntitiesChan
	pkg.Logrus(cons.INFO, "Found total data %d in postgres", len(usersEntities))

	if len(usersEntities) > 0 {
		cdcTimeUnix, err := helper.TimeStampToUnix(start_at)
		if err != nil {
			errChan <- err
			return
		}

		usersRepositorie := repo.NewUsersMeilisearchRepositorie(s.ctx, s.mls)

		filterAttributes := []string{"id", "deleted_at", "created_at", "updated_at"}
		sortAttributes := []string{"created_at"}

		if err := usersRepositorie.UpdateFilterableAttributes(filterAttributes...); err != nil {
			errChan <- err
			return
		}

		if err := usersRepositorie.UpdateSortableAttributes(sortAttributes...); err != nil {
			errChan <- err
			return
		}

		var insertDocFound, updateDocFound *bool
		usersDocEntitie := entitie.UsersDocument{}
		usersUpdateDocEntities := []entitie.UsersDocument{}
		usersInsertDocEntities := []entitie.UsersDocument{}

		for _, userEntity := range usersEntities {
			createdAtUnix, err := helper.TimeStampToUnix(userEntity.CreatedAt.Format(time.RFC3339))
			if err != nil {
				errChan <- err
				return
			}

			updatedAtUnix, err := helper.TimeStampToUnix(userEntity.UpdatedAt.Time.Format(time.RFC3339))
			if err != nil {
				errChan <- err
				return
			}

			deletedAtUnix, err := helper.TimeStampToUnix(userEntity.DeletedAt.Time.Format(time.RFC3339))
			if err != nil {
				errChan <- err
				return
			}

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
				errChan <- err
				return
			}

			createdAtFilter := fmt.Sprintf("updated_at = %d AND created_at > %d", defaultTimeUnix, cdcTimeUnix)
			updatedAtFilter := fmt.Sprintf("updated_at > %d", cdcTimeUnix)

			filter := fmt.Sprintf("deleted_at = %d AND id = '%s' AND (%s) OR (%s)", defaultTimeUnix, usersDocEntitie.ID, createdAtFilter, updatedAtFilter)
			fields := []string{
				"id",
				"name",
				"email",
				"phone",
				"date_of_birth",
				"age",
				"address",
				"city",
				"state",
				"direction",
				"country",
				"postal_code",
				"created_at",
				"updated_at",
				"deleted_at",
			}

			fetchSearch := meilisearch.DocumentsQuery{Filter: filter, Fields: fields}
			usersFetchDocuments, err := usersRepositorie.Find(&fetchSearch)

			if err != nil {
				errChan <- err
				return
			}

			if usersFetchDocuments.Results != nil {
				isTrue := true
				updateDocFound = &isTrue
				usersUpdateDocEntities = append(usersUpdateDocEntities, usersDocEntitie)
			} else {
				isTrue := true
				insertDocFound = &isTrue
				usersInsertDocEntities = append(usersInsertDocEntities, usersDocEntitie)
			}
		}

		if updateDocFound != nil && *updateDocFound {
			pkg.Logrus(cons.INFO, "Total data %d updated to meilisearch success", len(usersUpdateDocEntities))
			if err := usersRepositorie.BulkUpdate(usersUpdateDocEntities); err != nil {
				errChan <- err
				return
			}
			updateDocFound = nil
		}

		if insertDocFound != nil && *insertDocFound {
			pkg.Logrus(cons.INFO, "Total data %d inserted to meilisearch success", len(usersInsertDocEntities))
			if err := usersRepositorie.BulkInsert(usersInsertDocEntities); err != nil {
				errChan <- err
				return
			}
			insertDocFound = nil
		}
	}
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
		wg := new(sync.WaitGroup)

		usersEntitiesChan := make(chan []entitie.UsersEntitie, 1)
		errChan := make(chan error)

		wg.Add(2)
		go s.updateUsers(s.findAllUsers(wg, start_at, usersEntitiesChan, errChan))
		wg.Wait()

		close(usersEntitiesChan)
		close(errChan)

		select {
		case err := <-errChan:
			if err != nil {
				pkg.Logrus(cons.ERROR, err)
				return
			}
		}
	}
}

func (s searchScheduler) breakRun(rds inf.IRedis, handler func(rds inf.IRedis)) {
	key := "SCHEDULER:SEARCH:BREAK"
	value := 1
	sync := 10

	result, err := rds.IncrBy(key, value)
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	if result >= sync {
		breakTime := time.Duration(time.Second * 60)
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

		now := time.Now()

		key := "WORKER:SEARCH:CDC"
		value := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()-10, now.Second(), now.Nanosecond(), now.Location()).Format(time.RFC3339)

		if err := rds.Set(key, value); err != nil {
			pkg.Logrus(cons.ERROR, err)
			return
		}
	}

	if result <= sync {
		handler(rds)
	}
}

func (s searchScheduler) SearchRun() {
	cron := pkg.NewCron()

	crontime := cons.Every30Seconds
	now := time.Now().Format(cons.DATE_TIME_FORMAT)

	rds, err := pkg.NewRedis(s.ctx, s.rds)
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	sch, _, err := cron.Handler("search scheduler", crontime, func() {
		pkg.Logrus(cons.INFO, fmt.Sprintf("Search scheduler is running %s - and execute at %s", now, crontime))
		s.breakRun(rds, s.searchHandler)
	})

	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		return
	}

	sch.Start()
}
