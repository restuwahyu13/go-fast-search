package helper

import (
	"math"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"

	cons "github.com/restuwahyu13/go-fast-search/shared/constants"
	"github.com/restuwahyu13/go-fast-search/shared/dto"
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

func MeiliSearchPublisher[T any](amqp inf.IRabbitMQ, secret string, id any, data T, isBulk bool, action string) error {
	usersDocReq := dto.MeiliSearchDocuments[T]{}
	usersDocReq.ID = id
	usersDocReq.Doc = "users"
	usersDocReq.Data = any(data).(T)
	usersDocReq.IsBulk = isBulk
	usersDocReq.Action = action

	amqp_req := dto.Request[dto.RabbitOptions]{}

	amqp_req.Option.ExchangeName = cons.EXCHANGE_NAME_SEARCH
	amqp_req.Option.ExchangeType = cons.EXCHANGE_TYPE_DIRECT
	amqp_req.Option.QueueName = cons.QUEUE_NAME_SEARCH

	amqp_req.Option.Body = usersDocReq
	amqp_req.Option.Args = rabbitmq.Table{cons.X_RABBIT_SECRET: secret}

	if err := amqp.Publisher(amqp_req); err != nil {
		return err
	}

	return nil
}

func SleepBackoff(req dto.Request[dto.SleepBackoff]) {
	cmd := req.Config.Redis.IncrBy(req.Body.Ctx, req.Body.Key, int64(req.Body.Count))
	if cmd.Err() != nil {
		logrus.Error(cmd.Err())
		return
	}

	count := int(cmd.Val())

	if count >= req.Body.Retry {
		waitoff := time.Duration(math.Pow(float64(req.Body.Backoff), float64(count))) * time.Second
		time.Sleep(waitoff)

		cmd := req.Config.Redis.Set(req.Body.Ctx, req.Body.Key, req.Body.Count, 0)
		if cmd.Err() != nil {
			logrus.Error(cmd.Err())
			return
		}
	}

	if count <= req.Body.Retry {
		backoff := time.Duration(math.Pow(float64(req.Body.Backoff), float64(count))) * time.Second
		time.Sleep(backoff)
	}
}

func TimeStampToUnix(timestamp any) (int64, error) {
	var tparse time.Time

	if reflect.TypeOf(timestamp).Kind() == reflect.String {
		originalLayout := time.RFC3339

		tparse, err := time.Parse(originalLayout, timestamp.(string))
		if err != nil {
			return -1, err
		}

		return tparse.Unix(), nil
	}

	tparse, err := time.Parse(time.RFC3339, timestamp.(time.Time).String())
	if err != nil {
		return -1, err
	}

	return tparse.Unix(), nil
}
