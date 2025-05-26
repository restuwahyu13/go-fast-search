package inf

import "github.com/go-co-op/gocron/v2"

type ICron interface {
	Handler(name, crontime string, task func()) (gocron.Scheduler, gocron.Job, error)
}
