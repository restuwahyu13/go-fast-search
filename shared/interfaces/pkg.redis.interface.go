package inf

import (
	"time"
)

type IRedis interface {
	Set(key string, value any) error
	SetEx(key string, expiration time.Duration, value any) error
	Get(key string) ([]byte, error)
	Exists(key string) (int64, error)
	Del(key string) (int64, error)
	HSetEx(key string, expiration time.Duration, values ...any) error
	HGet(key string, field string) ([]byte, error)
	IncrBy(key string, value int) (int, error)
	TTL(key string) (int, error)
}
