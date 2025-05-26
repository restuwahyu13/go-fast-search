package pkg

import (
	"github.com/go-co-op/gocron/v2"

	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type cron struct{}

func NewCron() inf.ICron {
	return cron{}
}

func (p cron) Handler(name, crontime string, task func()) (gocron.Scheduler, gocron.Job, error) {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		return nil, nil, err
	}

	job, err := scheduler.NewJob(gocron.CronJob(crontime, true), gocron.NewTask(task), gocron.WithName(name))
	if err != nil {
		return nil, nil, err
	}

	return scheduler, job, nil
}
