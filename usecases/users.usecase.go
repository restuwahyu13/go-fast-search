package usecase

import (
	"context"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
)

type usersUsecase struct {
	service inf.IUsersService
}

func NewUsersUsecase(options dto.UsecaseOptions[inf.IUsersService]) inf.IUsersUsecase {
	return usersUsecase{service: options.SERVICE}
}

func (u usersUsecase) Ping(ctx context.Context) opt.Response {
	return u.service.Ping(ctx)
}

func (u usersUsecase) CreateUsers(ctx context.Context, req dto.Request[dto.CreateUsersDTO]) opt.Response {
	return u.service.CreateUsers(ctx, req)
}

func (u usersUsecase) UpdateUsers(ctx context.Context, req dto.Request[dto.UpdateUsersDTO]) opt.Response {
	return u.service.UpdateUsers(ctx, req)
}

func (u usersUsecase) FindAllUsers(ctx context.Context) opt.Response {
	return u.service.FindAllUsers(ctx)
}
