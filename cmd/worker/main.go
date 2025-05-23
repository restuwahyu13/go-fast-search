package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/wagslane/go-rabbitmq"

	config "github.com/restuwahyu13/go-fast-search/configs"
	con "github.com/restuwahyu13/go-fast-search/internal/infrastructure/connections"
	worker "github.com/restuwahyu13/go-fast-search/internal/infrastructure/workers"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type (
	IWorker interface {
		Listener()
	}

	Worker struct {
		CTX  context.Context
		ENV  dto.Request[dto.Environtment]
		DB   *bun.DB
		RDS  *redis.Client
		AMQP *rabbitmq.Conn
		MLS  meilisearch.ServiceManager
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

	req := dto.Request[Worker]{}
	req.Option = Worker{
		CTX:  ctx,
		ENV:  env,
		DB:   db,
		RDS:  rds,
		AMQP: amqp,
		MLS:  mls,
	}

	app := NewWorker(req)
	app.Listener()
}

func NewWorker(req dto.Request[Worker]) IWorker {
	return Worker{
		CTX:  req.Option.CTX,
		ENV:  req.Option.ENV,
		DB:   req.Option.DB,
		RDS:  req.Option.RDS,
		AMQP: req.Option.AMQP,
		MLS:  req.Option.MLS,
	}
}

func (i Worker) Register() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go worker.NewSearchWorker(dto.WorkerOptions{
		CTX:  i.CTX,
		ENV:  i.ENV,
		DB:   i.DB,
		RDS:  i.RDS,
		AMQP: i.AMQP,
		MLS:  i.MLS,
	}).SearchRun(wg)

	go worker.NewDeadLetterQueueWorker(dto.WorkerOptions{
		CTX:  i.CTX,
		ENV:  i.ENV,
		DB:   i.DB,
		RDS:  i.RDS,
		AMQP: i.AMQP,
		MLS:  i.MLS,
	}).DeadLetterQueueRun(wg)

	wg.Wait()
}

func (i Worker) Listener() {
	i.Register()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM, syscall.SIGABRT, syscall.SIGUSR1)

	for {
		select {
		case <-ch:
			if i.ENV.Config.APP.ENV != cons.DEV {
				time.Sleep(time.Second * 3)
			} else {
				time.Sleep(time.Second * 10)
			}

			os.Exit(0)

		default:
			time.Sleep(time.Second * 3)
		}
	}
}
