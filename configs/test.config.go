package config

import (
	"context"
	"errors"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/wagslane/go-rabbitmq"

	con "github.com/restuwahyu13/go-fast-search/internal/infrastructure/connections"
	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	helper "github.com/restuwahyu13/go-fast-search/shared/helpers"
	"github.com/restuwahyu13/go-fast-search/shared/pkg"
)

type (
	Test struct {
		CTX  context.Context
		ENV  dto.Request[dto.Environtment]
		DB   *bun.DB
		RDS  *redis.Client
		AMQP *rabbitmq.Conn
		MLS  meilisearch.ServiceManager
	}
)

var (
	err error
	env dto.Request[dto.Environtment]
)

func init() {
	transform := helper.NewTransform()

	env_res, err := NewEnvirontment(".env", ".", "env")
	if err != nil {
		pkg.Logrus(cons.FATAL, err)
		return
	}

	if env_res != nil {
		if err := transform.ResToReq(env_res, &env); err != nil {
			pkg.Logrus(cons.FATAL, err)
		}
	}
}

func NewTest() Test {
	ctx := context.Background()

	db, err := con.SqlConnection(ctx, env)
	if err != nil {
		pkg.Logrus(cons.FATAL, err)
	}
	defer db.Close()

	rds, err := con.RedisConnection(env)
	if err != nil {
		pkg.Logrus(cons.FATAL, err)
	}
	defer rds.Close()

	amqp, err := con.RabbitConnection(env)
	if err != nil {
		pkg.Logrus(cons.FATAL, err)
	}
	defer amqp.Close()

	mls := con.MeiliSearchConnection(env)
	if !mls.IsHealthy() {
		pkg.Logrus(cons.FATAL, errors.New("meilisearch is not healthy"))
	}
	defer mls.Close()

	return Test{
		CTX:  ctx,
		ENV:  env,
		DB:   db,
		RDS:  rds,
		AMQP: amqp,
		MLS:  mls,
	}
}
