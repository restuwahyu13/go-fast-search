package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type usersRoute struct {
	router     chi.Router
	controller inf.IUsersController
}

func NewUsersRoute(options dto.RouteOptions[inf.IUsersController]) {
	route := usersRoute{router: options.ROUTER, controller: options.CONTROLLER}

	route.router.Route(helper.Version("users"), func(r chi.Router) {
		r.Post("/", route.controller.CreateUsers)
		r.Get("/", route.controller.FindAllUsers)
		r.Put("/{id}", route.controller.UpdateUsers)
	})
}
