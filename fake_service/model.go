package fake_service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ypapax/status_check/queue"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Port struct {
	From        int   `yaml:"from"`
	To          int   `yaml:"to"`
	StatusCodes []int `yaml:"codes"`
	DelayMS     []int `yaml:"delay_ms"`
}

type Config struct {
	Ports []Port `yaml:"ports"`
}

func ParseConf(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal([]byte(b), &c); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &c, nil
}

func Serve(config *Config) error {
	serversCount := 0
	for _, ports := range config.Ports {
		for i := ports.From; i <= ports.To; i++ {
			serversCount++
			go serverOnPort(i, ports.StatusCodes, ports.DelayMS)
		}
	}
	logrus.Infof("%+v servers are listening", serversCount)
	forever := make(chan bool)
	<-forever
	return nil
}

func serverOnPort(port int, statusCodes []int, delayMS []int) error {
	router := mux.NewRouter().StrictSlash(true)
	var statusQueue = queue.New(statusCodes)
	var delayQueue = queue.New(delayMS)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if len(delayMS) > 0 {
			sleepDur := time.Duration(delayQueue.Next()) * time.Millisecond
			logrus.Tracef("sleeping for %s", sleepDur)
			time.Sleep(sleepDur)
		}
		w.WriteHeader(statusQueue.Next())
		if _, err := w.Write([]byte(fmt.Sprintf("%+v", statusCodes))); err != nil {
			logrus.Error(err)
		}
	}).Methods("GET")
	server := http.NewServeMux()
	server.Handle("/", router)
	bind := fmt.Sprintf("0.0.0.0:%d", port)
	logrus.Tracef("listening on %+v with status codes %+v", bind, statusCodes)
	if err := http.ListenAndServe(bind, server); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func AllPorts(conf *Config) []int {
	var allPorts []int
	for _, p := range conf.Ports {
		for i := p.From; i <= p.To; i++ {
			allPorts = append(allPorts, i)
		}
	}
	return allPorts
}
