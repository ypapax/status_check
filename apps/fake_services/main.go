package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check/fake_service"
)

func main() {
	rand.Seed(time.Now().Unix())
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)
	var confPath string
	flag.StringVar(&confPath, "conf", "fake-services.test.conf.yaml", "path to a config file")
	flag.Parse()
	f, err := os.Open(confPath)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	c, err := fake_service.ParseConf(f)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	if err := fake_service.Serve(c); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
