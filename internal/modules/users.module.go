package module

import (
	service "github.com/restuwahyu13/go-fast-search/domain/services"
	controller "github.com/restuwahyu13/go-fast-search/internal/adapters/http/controllers"
	route "github.com/restuwahyu13/go-fast-search/internal/adapters/http/routes"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	usecase "github.com/restuwahyu13/go-fast-search/usecases"
)

func NewUsersModule[IService any](options dto.ModuleOptions) {
	service := service.NewUsersService(dto.ServiceOptions{ENV: options.ENV, DB: options.DB, RDS: options.RDS})

	usecase := usecase.NewUsersUsecase(dto.UsecaseOptions[inf.IUsersService]{SERVICE: service})

	controller := controller.NewUsersController(dto.ControllerOptions[inf.IUsersUsecase]{USECASE: usecase})

	route.NewUsersRoute(dto.RouteOptions[inf.IUsersController]{ROUTER: options.ROUTER, CONTROLLER: controller})
}
