package helper

import (
	"reflect"
	"time"

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

func SleepBackoff[T any](req dto.Request[dto.SleepBackoff], handler func() (T, error)) (T, error) {
	cmdIncrBy := req.Config.Redis.IncrBy(req.Body.Ctx, req.Body.Key, req.Body.Count)
	if err := cmdIncrBy.Err(); err != nil {
		return any(nil).(T), err
	}

	if cmdIncrBy.Val() >= req.Body.Retry {
		breakTime := time.Duration(time.Second * time.Duration(req.Body.BackOffTime))

		cmdTTL := req.Config.Redis.TTL(req.Body.Ctx, req.Body.Key)
		if err := cmdTTL.Err(); err != nil {
			return any(nil).(T), err
		}

		if cmdTTL.Val() < 1 {
			cmdSet := req.Config.Redis.SetEx(req.Body.Ctx, req.Body.Key, cmdIncrBy.Val(), breakTime)
			if err := cmdSet.Err(); err != nil {
				return any(nil).(T), err
			}

		} else if cmdTTL.Val() > 1 && cmdTTL.Val() < 3 {
			cmdSet := req.Config.Redis.Set(req.Body.Ctx, req.Body.Key, req.Body.Count, 0)
			if err := cmdSet.Err(); err != nil {
				return any(nil).(T), err
			}
		}
	}

	if cmdIncrBy.Val() <= req.Body.Retry {
		return handler()
	}

	return any(nil).(T), nil
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
