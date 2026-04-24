package task

import (
	"github.com/robfig/cron/v3"
	"mochat-api-server/internal/pkg/logger"
)

var Scheduler *cron.Cron

func InitScheduler() {
	Scheduler = cron.New(cron.WithSeconds())
	logger.Sugar.Info("cron scheduler initialized")
}

func StartScheduler() {
	if Scheduler != nil {
		Scheduler.Start()
		logger.Sugar.Info("cron scheduler started")
	}
}

func StopScheduler() {
	if Scheduler != nil {
		ctx := Scheduler.Stop()
		<-ctx.Done()
		logger.Sugar.Info("cron scheduler stopped")
	}
}

func AddFunc(spec string, cmd func()) error {
	_, err := Scheduler.AddFunc(spec, cmd)
	if err != nil {
		return err
	}
	logger.Sugar.Infof("registered cron task: %s", spec)
	return nil
}
