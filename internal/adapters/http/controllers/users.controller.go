package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	gpc "github.com/restuwahyu13/go-playground-converter"

	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type usersController struct {
	usecase inf.IUsersUsecase
}

func NewUsersController(options dto.ControllerOptions[inf.IUsersUsecase]) inf.IUsersController {
	return usersController{usecase: options.USECASE}
}

func (c usersController) Ping(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := opt.Response{}

	if res = c.usecase.Ping(ctx); res.StatCode >= http.StatusBadRequest {
		if res.StatCode >= http.StatusInternalServerError {
			pkg.Logrus(cons.ERROR, res.ErrMsg)
			res.ErrMsg = cons.DEFAULT_ERR_MSG
		}

		helper.Api(rw, r, res)
		return
	}

	helper.Api(rw, r, res)
	return
}

func (c usersController) CreateUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parser := helper.NewParser()

	res := opt.Response{}
	req := dto.Request[dto.CreateUsersDTO]{}

	if err := parser.Decode(r.Body, &req.Body); err != nil {
		pkg.Logrus(cons.ERROR, err)
		helper.Api(rw, r, res)
		return
	}

	errors, err := gpc.Validator(req.Body)
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		helper.Api(rw, r, res)
		return
	}

	if errors != nil {
		res.StatCode = http.StatusUnprocessableEntity
		res.Errors = errors.Errors

		helper.Api(rw, r, res)
		return
	}

	if res = c.usecase.CreateUsers(ctx, req); res.StatCode >= http.StatusBadRequest {
		if res.StatCode >= http.StatusInternalServerError {
			pkg.Logrus(cons.ERROR, res.ErrMsg)
			res.ErrMsg = cons.DEFAULT_ERR_MSG
		}

		helper.Api(rw, r, res)
		return
	}

	helper.Api(rw, r, res)
	return
}

func (c usersController) UpdateUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parser := helper.NewParser()

	res := opt.Response{}
	req := dto.Request[dto.UpdateUsersDTO]{}

	req.Body.ID = chi.URLParam(r, "id")

	if err := parser.Decode(r.Body, &req.Body); err != nil {
		pkg.Logrus(cons.ERROR, err)
		helper.Api(rw, r, res)
		return
	}

	errors, err := gpc.Validator(req.Body)
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		helper.Api(rw, r, res)
		return
	}

	if errors != nil {
		res.StatCode = http.StatusUnprocessableEntity
		res.Errors = errors.Errors

		helper.Api(rw, r, res)
		return
	}

	if res = c.usecase.UpdateUsers(ctx, req); res.StatCode >= http.StatusBadRequest {
		if res.StatCode >= http.StatusInternalServerError {
			pkg.Logrus(cons.ERROR, res.ErrMsg)
			res.ErrMsg = cons.DEFAULT_ERR_MSG
		}

		helper.Api(rw, r, res)
		return
	}

	helper.Api(rw, r, res)
	return
}

func (c usersController) FindAllUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	transform := helper.NewTransform()

	res := opt.Response{}
	req := dto.Request[dto.MeiliSearchDocumentsQuery]{}

	if err := transform.QueryToStruct(r.URL.Query().Encode(), &req.Query); err != nil {
		pkg.Logrus(cons.ERROR, err)
		helper.Api(rw, r, res)
		return
	}

	errors, err := gpc.Validator(req.Query)
	if err != nil {
		pkg.Logrus(cons.ERROR, err)
		helper.Api(rw, r, res)
		return
	}

	if errors != nil {
		res.StatCode = http.StatusUnprocessableEntity
		res.Errors = errors.Errors

		helper.Api(rw, r, res)
		return
	}

	if res = c.usecase.FindAllUsers(ctx, req); res.StatCode >= http.StatusBadRequest {
		if res.StatCode >= http.StatusInternalServerError {
			pkg.Logrus(cons.ERROR, res.ErrMsg)
			res.ErrMsg = cons.DEFAULT_ERR_MSG
		}

		helper.Api(rw, r, res)
		return
	}

	res.Data = req.Query

	helper.Api(rw, r, res)
	return
}
