package checker

import (
	"context"
	"fmt"
	"time"

	"github.com/ypapax/status_check"
	"github.com/ypapax/status_check/config"
	"github.com/ypapax/status_check/job"
	web_service "github.com/ypapax/status_check/web-service"

	"github.com/sirupsen/logrus"
)

func Check(conf config.Config) error {
	if err := conf.Validate(); err != nil {
		return err
	}

	webServicesService, _, err := status_check.DbServices(conf.DbType, conf.ConnectionString)
	if err != nil {
		logrus.Error(err)
		return err
	}

	statusPipe, err := status_check.PipeServices(conf.PipeType, conf.Kafka.StatusTopic, conf.Kafka.ClientID, conf.Kafka.Brokers)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if err := createServices(webServicesService, conf.Addresses); err != nil {
		logrus.Error(err)
		return err
	}

	ss, err := webServicesService.FindAllWebServices()
	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Println("DbServices ", len(ss))

	parent := context.Background()
	var jobs = make(chan job.Job, conf.Workers)
	go func() {
		for _, s := range ss {
			j := job.Job{Service: s}
			logrus.Printf("sending job %+v to jobs channel", j)
			jobs <- j
			logrus.Println("jobs size", len(jobs))
		}
	}()
	for i := 0; i < conf.Workers; i++ {
		logrus.Println("starting worker ", i)
		go func(worker int) {
			for {
				func() {
					logrus.Trace("before getting job item to process")
					j := <-jobs
					logrus.Tracef("worker %+v receives job %+v", worker, j)
					defer func() {
						go func() {
							logrus.Tracef("about to write job %+v to jobs chan", j)
							time.Sleep(conf.CheckPeriod)
							jobs <- j
							logrus.Tracef("have written job %+v to jobs chan", j)
						}()
					}()
					if !j.LastCheckedTime.IsZero() && time.Since(j.LastCheckedTime) < conf.CheckPeriod {
						return
					}
					j.LastCheckedTime = time.Now()

					status, err := webServicesService.CheckStatus(&j.Service, conf.Schemas)
					if err != nil {
						err := fmt.Errorf("error %+v for checking service %+v", err, j.Service)
						logrus.Error(err)
						return
					}
					logrus.Tracef("publishing status to a pipe %+v", *status)
					if err := statusPipe.Publish(parent, *status); err != nil {
						err := fmt.Errorf("error %+v for checking service %+v", err, j.Service)
						logrus.Error(err)
						return
					}
				}()
			}
		}(i)
	}

	forever := make(chan bool)
	<-forever
	return nil
}

func createServices(webServicesService web_service.Service, addresses []string) error {
	var ss = make([]web_service.WebService, len(addresses))
	for i, a := range addresses {
		ss[i] = web_service.WebService{Address: a}
	}
	if err := webServicesService.CreateWebServices(ss); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
