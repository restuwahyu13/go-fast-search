package main

import (
	"compress/zlib"
	"context"
	"errors"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/meilisearch/meilisearch-go"
	"github.com/oxequa/grace"
	"github.com/redis/go-redis/v9"
	"github.com/unrolled/secure"
	"github.com/uptrace/bun"
	"github.com/wagslane/go-rabbitmq"

	config "github.com/restuwahyu13/go-fast-search/configs"
	con "github.com/restuwahyu13/go-fast-search/internal/infrastructure/connections"
	module "github.com/restuwahyu13/go-fast-search/internal/modules"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type (
	IApi interface {
		Middleware()
		Module()
		Listener()
	}

	Api struct {
		ENV     dto.Request[dto.Environtment]
		ENV_RES *opt.Environtment
		ROUTER  *chi.Mux
		DB      *bun.DB
		RDS     *redis.Client
		AMQP    *rabbitmq.Conn
		MLS     meilisearch.ServiceManager
	}
)

var (
	err     error
	env     dto.Request[dto.Environtment]
	env_res *opt.Environtment
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	transform := helper.NewTransform()

	env_res, err = config.NewEnvirontment(".env", ".", "env")
	if err != nil {
		pkg.Logrus(cons.FATAL, err)
		return
	}

	if env_res != nil {
		if err := transform.ResToReq(env_res, &env.Config); err != nil {
			pkg.Logrus(cons.FATAL, err)
			return
		}
	}
}

func main() {
	ctx := context.Background()
	router := chi.NewRouter()

	db, err := con.SqlConnection(ctx, env)
	if err != nil {
		pkg.Logrus(cons.FATAL, err)
		return
	}
	defer db.Close()

	rds, err := con.RedisConnection(env)
	if err != nil {
		pkg.Logrus(cons.FATAL, err)
		return
	}
	defer rds.Close()

	amqp, err := con.RabbitConnection(env)
	if err != nil {
		pkg.Logrus(cons.FATAL, err)
		return
	}
	defer amqp.Close()

	mls := con.MeiliSearchConnection(env)
	if !mls.IsHealthy() {
		pkg.Logrus(cons.FATAL, errors.New("meilisearch is not healthy"))
		return
	}
	defer mls.Close()

	req := dto.Request[Api]{}
	req.Option = Api{
		ENV:     env,
		ENV_RES: env_res,
		ROUTER:  router,
		DB:      db,
		RDS:     rds,
		AMQP:    amqp,
		MLS:     mls,
	}

	app := NewApi(req)
	app.Middleware()
	app.Module()
	app.Listener()
}

func NewApi(req dto.Request[Api]) IApi {
	return Api{
		ENV:     req.Option.ENV,
		ENV_RES: req.Option.ENV_RES,
		ROUTER:  req.Option.ROUTER,
		DB:      req.Option.DB,
		RDS:     req.Option.RDS,
		AMQP:    req.Option.AMQP,
		MLS:     req.Option.MLS,
	}
}

func (a Api) Middleware() {
	if a.ENV.Config.APP.ENV != cons.PROD {
		a.ROUTER.Use(middleware.Logger)
	}

	a.ROUTER.Use(middleware.Recoverer)
	a.ROUTER.Use(middleware.RealIP)
	a.ROUTER.Use(middleware.NoCache)
	a.ROUTER.Use(middleware.GetHead)
	a.ROUTER.Use(middleware.Compress(zlib.BestCompression))
	a.ROUTER.Use(middleware.AllowContentType("application/json"))
	a.ROUTER.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             900,
	}))
	a.ROUTER.Use(secure.New(secure.Options{
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		STSSeconds:           900,
	}).Handler)

	a.ROUTER.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helper.Api(w, r, opt.Response{
			StatCode: http.StatusMethodNotAllowed,
			ErrMsg:   "Route method not allowed",
		})
	}))

	a.ROUTER.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helper.Api(w, r, opt.Response{
			StatCode: http.StatusNotFound,
			ErrMsg:   "Route not found",
		})
	}))
}

func (a Api) Module() {
	module.NewUsersModule[inf.IUsersService](dto.ModuleOptions{
		ENV:    a.ENV,
		DB:     a.DB,
		RDS:    a.RDS,
		AMQP:   a.AMQP,
		MLS:    a.MLS,
		ROUTER: a.ROUTER,
	})
}

func (a Api) Listener() {
	err := pkg.Graceful(env, func() opt.Graceful {
		return opt.Graceful{HANDLER: a.ROUTER, ENV: a.ENV_RES}
	})

	recover := grace.Recover(&err)
	recover.Stack()

	if err != nil {
		pkg.Logrus(cons.FATAL, err)
		os.Exit(1)
		return
	}
}
