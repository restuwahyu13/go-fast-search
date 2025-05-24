package repo

import (
	"context"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"

	entitie "github.com/restuwahyu13/go-fast-search/domain/entities"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
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

func (r usersMeilisearchRepositorie) Search(query string, attributes []string, filter *meilisearch.SearchRequest) (*opt.UsersSearch, error) {
	parser := helper.NewParser()

	result := new(meilisearch.SearchResponse)
	res := new(opt.UsersSearch)

	if err := r.meilisearch.CreateCollection("users", "id", r.doc); err != nil {
		return nil, err
	}

	if attributes != nil {
		if _, err := r.meilisearch.CreateFilterableAttributes("users", attributes); err != nil {
			return nil, err
		}
	}

	err := r.meilisearch.Like("users", query, filter, result)
	if err != nil {
		return nil, err
	}

	resultByte, err := parser.Marshal(result)
	if err != nil {
		return nil, err
	}

	if err := parser.Unmarshal(resultByte, res); err != nil {
		return nil, err
	}

	return res, nil
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
