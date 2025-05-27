package repo

import (
	"context"
	"time"

	"github.com/guregu/null/v6/zero"
	"github.com/uptrace/bun"

	entitie "github.com/restuwahyu13/go-fast-search/domain/entities"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type usersRepositorie struct {
	ctx     context.Context
	db      *bun.DB
	entitie *entitie.UsersEntitie
}

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
