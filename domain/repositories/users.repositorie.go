package repo

import (
	"context"

	"github.com/uptrace/bun"

	entitie "github.com/restuwahyu13/go-fast-search/domain/entities"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type usersRepositorie struct {
	ctx     context.Context
	db      *bun.DB
	entitie *entitie.UsersEntitie
}

func NewUsersRepository(ctx context.Context, db *bun.DB) inf.IUsersRepositorie {
	return usersRepositorie{ctx: ctx, db: db, entitie: new(entitie.UsersEntitie)}
}

func (r usersRepositorie) Find(dest any) *bun.SelectQuery {
	users := []entitie.UsersEntitie{}

	if dest != nil {
		return r.db.NewSelect().Model(dest)
	}

	return r.db.NewSelect().Model(&users)
}

func (r usersRepositorie) FindOne(dest any) *bun.SelectQuery {
	if dest != nil {
		return r.db.NewSelect().Model(dest)
	}

	return r.db.NewSelect().Model(r.entitie)
}

func (r usersRepositorie) Create(dest any) *bun.InsertQuery {
	if dest != nil {
		return r.db.NewInsert().Model(dest)
	}

	return r.db.NewInsert().Model(r.entitie)
}

func (r usersRepositorie) Update(dest any) *bun.UpdateQuery {
	if dest != nil {
		return r.db.NewUpdate().Model(dest)
	}

	return r.db.NewUpdate().Model(r.entitie)
}

func (r usersRepositorie) Delete(dest any) *bun.DeleteQuery {
	if dest != nil {
		return r.db.NewDelete().Model(dest)
	}

	return r.db.NewDelete().Model(r.entitie)
}
