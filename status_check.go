package status_check

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check/database/psql"
	"github.com/ypapax/status_check/job"
	"github.com/ypapax/status_check/status"
	"github.com/ypapax/status_check/statuses"
	web_service "github.com/ypapax/status_check/web-service"
	"gopkg.in/yaml.v2"
)

const writeStatusesToDbPeriod = 5 * time.Second

type Config struct {
	Bind             string        `yaml:"bind"`
	CheckPeriod      time.Duration `yaml:"check_period"`
	DbType           string        `yaml:"db_type"`
	ConnectionString string        `yaml:"connection_string"`
	Workers          int           `yaml:"workers"`
	Schemas          []string      `yaml:"schemas"`
	Addresses        []string      `yaml:"addresses"`
}

func Serve(conf Config) error {
	if err := validateConf(conf); err != nil {
		return err
	}
	var webServicesService web_service.Service
	var statusService status.Service
	switch conf.DbType {
	case "psql":
		db, err := psql.GetConnection(conf.ConnectionString)
		if err != nil {
			logrus.Error(err)
			return err
		}
		defer db.Close()
		serviceRepo := psql.NewPostgresServiceRepository(db)
		statusRepo := psql.NewPostgresStatusRepository(db)
		webServicesService = web_service.NewService(serviceRepo)
		statusService = status.NewService(statusRepo)
	default:
		err := fmt.Errorf("db type '%+v' is not supported", conf.DbType)
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
	for _, s := range ss {
		jobs <- job.Job{Service: &s}
	}
	for i := 0; i < conf.Workers; i++ {
		go func() {
			j := <-jobs
			defer func() {
				jobs <- j
			}()
			if !j.LastCheckedTime.IsZero() && time.Since(j.LastCheckedTime) < conf.CheckPeriod {
				return
			}
			j.LastCheckedTime = time.Now()

			status, err := webServicesService.CheckStatus(j.Service, conf.Schemas)
			if err != nil {
				err := fmt.Errorf("error %+v for checking service %+v", err, j.Service)
				logrus.Error(err)
				return
			}
			sts.Add(*status)
		}()
	}
	forever := make(chan bool)
	<-forever

	return nil
}

func ParseConf(reader io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal([]byte(b), &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func validateConf(c Config) error {
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
