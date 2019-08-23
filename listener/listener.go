package listener

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check"
	"github.com/ypapax/status_check/config"
	"github.com/ypapax/status_check/status"
	"github.com/ypapax/status_check/statuses"
)

const writeStatusesToDbPeriod = 5 * time.Second

func ListenStatus(conf *config.Config) error {
	_, statusService, err := status_check.DbServices(conf.DbType, conf.ConnectionString)
	if err != nil {
		logrus.Error(err)
		return err
	}

	statusPipe, err := status_check.PipeServices(conf.PipeType, conf.Kafka.StatusTopic, conf.Kafka.ClientID, conf.Kafka.Brokers)
	if err != nil {
		logrus.Error(err)
		return err
	}
	parent := context.Background()
	var statusChan = make(chan status.Status)
	var errs = make(chan error)
	logrus.Printf("statusPipe: %+v", statusPipe)
	go statusPipe.Listen(parent, statusChan, errs)
	var sts statuses.Statuses
	go func() {
		t := time.NewTicker(writeStatusesToDbPeriod)
		for {
			<-t.C
			ss := sts.GetAll()
			if len(ss) == 0 {
				continue
			}
			if err := statusService.CreateStatus(ss); err != nil {
				logrus.Error(err)
			}
		}
	}()

	for {
		logrus.Tracef("before selecting a channel")
		select {
		case st := <-statusChan:
			sts.Add(st)
			logrus.Tracef("received status %+v", st)
		case err := <-errs:
			logrus.Error(err)
		}
	}
}
