package repo

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"

	entitie "github.com/restuwahyu13/go-fast-search/domain/entities"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type (
	usersRepositorie struct {
		ctx     context.Context
		db      *bun.DB
		entitie *entitie.UsersEntitie
	}

	usersMeilisearchRepositorie struct {
		ctx         context.Context
		meilisearch inf.IMeiliSearch
		doc         *entitie.UsersDocument
	}
)

/*
* ===============================================
*  REPOSITORIE FOR RDBMS
* ===============================================
 */

func NewUsersRepositorie(ctx context.Context, db *bun.DB) inf.IUsersRepositorie {
	return usersRepositorie{ctx: ctx, db: db, entitie: new(entitie.UsersEntitie)}
}

func (r usersRepositorie) Find() *bun.SelectQuery {
	return r.db.NewSelect().Model(r.entitie)
}

func (r usersRepositorie) FindOne() *bun.SelectQuery {
	return r.db.NewSelect().Model(r.entitie)
}

func (r usersRepositorie) Insert(entitie entitie.UsersEntitie, column string, dest ...any) error {
	sqlb := r.db.NewInsert().Model(&entitie)

	if column != "" && dest != nil {
		result, err := sqlb.Returning(column).Exec(r.ctx, dest...)
		if err != nil {
			return nil

		} else if rows, err := result.RowsAffected(); err != nil || rows < 1 {
			if err != nil {
				return err
			}

			return cons.NO_ROWS_AFFECTED

		}
	} else {
		result, err := sqlb.Exec(r.ctx)
		if err != nil {
			return err

		} else if rows, err := result.RowsAffected(); err != nil || rows < 1 {
			if err != nil {
				return err
			}

			return cons.NO_ROWS_AFFECTED

		}
	}

	return nil
}

func (r usersRepositorie) Update(entitie entitie.UsersEntitie, column string, dest ...any) error {
	sqlb := r.db.NewUpdate().Model(&entitie)

	if column != "" && dest != nil {
		result, err := sqlb.Returning(column).Where("deleted_at IS NULL AND id = ?", entitie.ID).OmitZero().Exec(r.ctx, dest...)
		if err != nil {
			return nil

		} else if rows, err := result.RowsAffected(); err != nil || rows < 1 {
			if err != nil {
				return err
			}

			return cons.NO_ROWS_AFFECTED

		}
	} else {
		result, err := sqlb.Where("deleted_at IS NULL AND id = ?", entitie.ID).OmitZero().Exec(r.ctx)
		if err != nil {
			return nil

		} else if rows, err := result.RowsAffected(); err != nil || rows < 1 {
			if err != nil {
				return err
			}

			return cons.NO_ROWS_AFFECTED

		}
	}

	return nil
}

func (r usersRepositorie) Delete(id string, dest any) error {
	r.entitie.DeletedAt = zero.TimeFrom(time.Now())
	r.entitie.UpdatedAt = zero.TimeFrom(time.Now())

	sqlb := r.db.NewUpdate().Model(r.entitie).Where("id = ?", id)

	if dest != nil {
		result, err := sqlb.Returning("*", dest).OmitZero().Exec(r.ctx)
		if err != nil {
			return nil

		} else if rows, err := result.RowsAffected(); err != nil || rows < 1 {
			if err != nil {
				return err
			}

			return cons.NO_ROWS_AFFECTED

		}
	}

	result, err := sqlb.OmitZero().Exec(r.ctx)
	if err != nil {
		return nil

	} else if rows, err := result.RowsAffected(); err != nil || rows < 1 {
		if err != nil {
			return err
		}

		return cons.NO_ROWS_AFFECTED

	}

	return nil
}

/*
* ===============================================
*  REPOSITORIE FOR MEILISEARCH
* ===============================================
 */

func NewUsersMeilisearchRepositorie(ctx context.Context, db meilisearch.ServiceManager) inf.IUsersMeiliSearchRepositorie {
	meilisearch := pkg.NewMeiliSearch(ctx, db)
	return usersMeilisearchRepositorie{ctx: ctx, meilisearch: meilisearch, doc: new(entitie.UsersDocument)}
}

func (r usersMeilisearchRepositorie) Search(query string, filter *meilisearch.SearchRequest) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error) {
	transform := helper.NewTransform()

	docResult := new(meilisearch.SearchResponse)
	docResultReformat := new(opt.MeiliSearchDocuments[[]entitie.UsersDocument])

	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return nil, err
	}

	err := r.meilisearch.Like("users", query, filter, docResult)
	if err != nil {
		return nil, err
	}

	if err := transform.SrcToDest(docResult, docResultReformat); err != nil {
		return nil, err
	}

	return docResultReformat, nil
}

func (r usersMeilisearchRepositorie) Find(filter *meilisearch.DocumentsQuery) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error) {
	transform := helper.NewTransform()

	docResult := new(meilisearch.DocumentsResult)
	docResultReformat := new(opt.MeiliSearchDocuments[[]entitie.UsersDocument])

	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return nil, err
	}

	err := r.meilisearch.Find("users", filter, docResult)
	if err != nil {
		return nil, err
	}

	if err := transform.SrcToDest(docResult, docResultReformat); err != nil {
		return nil, err
	}

	return docResultReformat, nil
}

func (r usersMeilisearchRepositorie) FindOne(id string, filter *meilisearch.DocumentQuery) (*entitie.UsersDocument, error) {
	res := new(entitie.UsersDocument)

	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return nil, err
	}

	err := r.meilisearch.FindOne("users", id, filter, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r usersMeilisearchRepositorie) Insert(value any) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.Insert("users", value); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) Update(id string, value any) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.Update("users", id, value); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) Delete(id string) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.Delete("users", id); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) BulkInsert(value any) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.BulkInsert("users", value); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) BulkUpdate(value any) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.BulkUpdate("users", value); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) BulkDelete(ids ...string) error {
	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return err
	}

	if _, err := r.meilisearch.BulkDelete("users", ids...); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) UpdateFilterableAttributes(attributes ...string) error {
	if _, err := r.meilisearch.UpdateFilterableAttributes("users", attributes); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) UpdateSearchableAttributes(attributes ...string) error {
	if _, err := r.meilisearch.UpdateSearchableAttributes("users", attributes); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) UpdateSortableAttributes(attributes ...string) error {
	if _, err := r.meilisearch.UpdateSortableAttributes("users", attributes); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) UpdateDisplayedAttributes(attributes ...string) error {
	if _, err := r.meilisearch.UpdateDisplayedAttributes("users", attributes); err != nil {
		return err
	}

	return nil
}

func (r usersMeilisearchRepositorie) ListUsersDocuments(req dto.Request[dto.MeiliSearchDocumentsQuery]) (*opt.MeiliSearchDocuments[[]entitie.UsersDocument], error) {
	transform := helper.NewTransform()

	deletedAtUnix, err := helper.TimeStampToUnix(time.Time{}.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("deleted_at = %d", deletedAtUnix)
	mlsReq := new(meilisearch.SearchRequest)

	mlsReq.Limit = req.Query.Limit
	mlsReq.Page = req.Query.Page
	mlsReq.Sort = []string{req.Query.SortBy + ":" + req.Query.Sort}

	if req.Query.MatchingStrategy != "" {
		switch req.Query.MatchingStrategy {

		case string(meilisearch.Last):
			mlsReq.MatchingStrategy = meilisearch.Last
			break

		case string(meilisearch.All):
			mlsReq.MatchingStrategy = meilisearch.All
			break

		case string(meilisearch.Frequency):
			mlsReq.MatchingStrategy = meilisearch.Frequency
			break
		}
	}

	if req.Query.HighlightAttributes != "" {
		mlsReq.AttributesToHighlight = strings.Split(req.Query.HighlightAttributes, ",")
	}

	if req.Query.Filter != nil && req.Query.FilterBy != "" {
		filterBy := strings.Split(req.Query.FilterBy, ",")

		if err := r.UpdateFilterableAttributes(filterBy...); err != nil {
			return nil, err
		}
	}

	if req.Query.Sort != "" && req.Query.SortBy != "" {
		sortBy := strings.Split(req.Query.SortBy, ",")

		if err := r.UpdateSortableAttributes(sortBy...); err != nil {
			return nil, err
		}

		mlsReq.Sort = []string{req.Query.SortBy + ":" + req.Query.Sort}
	}

	usersFilterDoc := new(dto.ListUsersFilterDTO)
	if err := transform.ReqToRes(&req.Query.Filter, usersFilterDoc); err != nil {
		return nil, err
	}

	if usersFilterDoc.StartDate != "" && usersFilterDoc.EndDate != "" {
		startDate, err := helper.TimeStampToUnix(usersFilterDoc.StartDate)
		if err != nil {
			return nil, err
		}

		endDate, err := helper.TimeStampToUnix(usersFilterDoc.EndDate)
		if err != nil {
			return nil, err
		}

		filter += fmt.Sprintf(" AND created_at > %d AND created_at < %d", startDate, endDate)
	}

	if usersFilterDoc.Age != "" {
		filter += fmt.Sprintf(" AND age = %s", usersFilterDoc.Age)
	}

	if usersFilterDoc.City != "" {
		filter += fmt.Sprintf(" AND age = %s", usersFilterDoc.City)
	}

	if usersFilterDoc.State != "" {
		filter += fmt.Sprintf(" AND state = %s", usersFilterDoc.State)
	}

	if usersFilterDoc.Direction != "" {
		filter += fmt.Sprintf(" AND direction = %s", usersFilterDoc.Direction)
	}

	if usersFilterDoc.Country != "" {
		filter += fmt.Sprintf(" AND country = %s", usersFilterDoc.Country)
	}

	resultUsersStats, err := r.meilisearch.GetStats("users")
	if err != nil {
		return nil, err
	}

	mlsReq.Filter = filter
	resultUsersDocuments, err := r.Search(req.Query.Search, mlsReq)
	if err != nil {
		return nil, err
	}

	resultUsersDocuments.Page = req.Query.Page + 1
	resultUsersDocuments.TotalPages = int64(math.Ceil(float64(resultUsersStats.NumberOfDocuments) / float64(req.Query.Limit)))
	resultUsersDocuments.Total = resultUsersStats.NumberOfDocuments

	return resultUsersDocuments, nil
}
