package status_check

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/ypapax/status_check/config"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check/database/psql"
	"github.com/ypapax/status_check/job"
	"github.com/ypapax/status_check/status"
	"github.com/ypapax/status_check/statuses"
	web_service "github.com/ypapax/status_check/web-service"
	"gopkg.in/yaml.v2"
)

const writeStatusesToDbPeriod = 5 * time.Second

func CheckServices(conf config.Config) error {
	if err := validateConf(conf); err != nil {
		return err
	}

	webServicesService, statusService, err := Services(conf.DbType, conf.ConnectionString)
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

	logrus.Println("Services ", len(ss))

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
					sts.Add(*status)
				}()
			}
		}(i)
	}

	forever := make(chan bool)
	<-forever
	return nil
}

func Services(dbType, connString string) (web_service.Service, status.Service, error) {
	var webServicesService web_service.Service
	var statusService status.Service

	switch dbType {
	case "psql":
		db, err := psql.GetConnection(connString, 10*time.Second)
		if err != nil {
			logrus.Error(err)
			return webServicesService, statusService, err
		}
		serviceRepo := psql.NewPostgresServiceRepository(db)
		statusRepo := psql.NewPostgresStatusRepository(db)
		webServicesService = web_service.NewService(serviceRepo)
		statusService = status.NewService(statusRepo)
	default:
		err := fmt.Errorf("db type '%+v' is not supported", connString)
		return webServicesService, statusService, err
	}
	return webServicesService, statusService, nil
}

func ParseConf(reader io.Reader) (*config.Config, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var c config.Config
	if err := yaml.Unmarshal([]byte(b), &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func validateConf(c config.Config) error {
	if len(c.Addresses) == 0 {
		return fmt.Errorf("Missing addresses list")
	}
	if len(c.Schemas) == 0 {
		return fmt.Errorf("Missing schemas list")
	}
	if c.Workers <= 0 {
		return fmt.Errorf("Workers amount should be positive")
	}
	if c.CheckPeriod == 0 {
		return fmt.Errorf("Empty check period")
	}
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
