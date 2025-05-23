package main

import (
	"compress/zlib"
	"context"
	"errors"
	"os"

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
		MLS:     req.Option.MLS,
	}
}

func (i Api) Middleware() {
	if i.ENV.Config.APP.ENV != cons.PROD {
		i.ROUTER.Use(middleware.Logger)
	}

	i.ROUTER.Use(middleware.Recoverer)
	i.ROUTER.Use(middleware.RealIP)
	i.ROUTER.Use(middleware.NoCache)
	i.ROUTER.Use(middleware.GetHead)
	i.ROUTER.Use(middleware.Compress(zlib.BestCompression))
	i.ROUTER.Use(middleware.AllowContentType("application/json"))
	i.ROUTER.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             900,
	}))
	i.ROUTER.Use(secure.New(secure.Options{
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		STSSeconds:           900,
	}).Handler)
}

func (i Api) Module() {
	module.NewUsersModule[inf.IUsersService](dto.ModuleOptions{
		ENV:    i.ENV,
		DB:     i.DB,
		RDS:    i.RDS,
		MLS:    i.MLS,
		ROUTER: i.ROUTER,
	})
}

func (i Api) Listener() {

	err := pkg.Graceful(env, func() opt.Graceful {
		return opt.Graceful{HANDLER: i.ROUTER, ENV: i.ENV_RES}
	})

	recover := grace.Recover(&err)
	recover.Stack()

	if err != nil {
		pkg.Logrus(cons.FATAL, err)
		os.Exit(1)
		return
	}
}
