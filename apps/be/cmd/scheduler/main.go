package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"

	config "github.com/restuwahyu13/go-fast-search/configs"
	con "github.com/restuwahyu13/go-fast-search/internal/infrastructure/connections"
	scheduler "github.com/restuwahyu13/go-fast-search/internal/infrastructure/schedulers"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	opt "github.com/restuwahyu13/go-fast-search/shared/output"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type (
	IScheduler interface {
		Listener()
	}

	Scheduler struct {
		CTX context.Context
		ENV dto.Request[dto.Environtment]
		DB  *bun.DB
		RDS *redis.Client
		MLS meilisearch.ServiceManager
	}

	syncOnce struct {
		searchScheduler *sync.Once
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

	mls := con.MeiliSearchConnection(env)
	if !mls.IsHealthy() {
		pkg.Logrus(cons.FATAL, errors.New("meilisearch is not healthy"))
		return
	}
	defer mls.Close()

	req := dto.Request[Scheduler]{}
	req.Option = Scheduler{
		CTX: ctx,
		ENV: env,
		DB:  db,
		RDS: rds,
		MLS: mls,
	}

	app := NewScheduler(req)
	app.Listener()
}

func NewScheduler(req dto.Request[Scheduler]) IScheduler {
	return Scheduler{
		CTX: req.Option.CTX,
		ENV: req.Option.ENV,
		DB:  req.Option.DB,
		RDS: req.Option.RDS,
		MLS: req.Option.MLS,
	}
}

func (w Scheduler) scheduler(wg *sync.WaitGroup, rso syncOnce) {
	defer wg.Done()

	rso.searchScheduler.Do(func() {
		scheduler.NewSearchScheduler(dto.SchedulerOptions{
			CTX: w.CTX,
			ENV: w.ENV,
			DB:  w.DB,
			RDS: w.RDS,
			MLS: w.MLS,
		}).SearchRun()
	})
}

func (w Scheduler) register(wg *sync.WaitGroup) {
	worker := runtime.NumCPU()
	searchSchedulerOnce := new(sync.Once)

	rso := syncOnce{
		searchScheduler: searchSchedulerOnce,
	}

	for i := 1; i <= worker; i++ {
		wg.Add(1)
		go w.scheduler(wg, rso)
	}
}

func (w Scheduler) Listener() {
	wg := sync.WaitGroup{}
	w.register(&wg)

	ctx, cancel := context.WithCancel(w.CTX)
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM, syscall.SIGABRT, syscall.SIGUSR1)

	for {
		select {
		case <-ch:
			if w.ENV.Config.APP.ENV != cons.DEV {
				time.Sleep(time.Second * 10)
			} else {
				time.Sleep(time.Second * 15)
			}

			wg.Wait()
			os.Exit(0)
			return

		case <-ctx.Done():
			wg.Wait()
			return

		default:
			wg.Wait()
			time.Sleep(time.Second * 5)
			pkg.Logrus(cons.INFO, "Scheduler is running")
		}
	}
}
