package main

import (
	"flag"
	"github.com/ypapax/logrus_conf"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check/fake_config"
	"github.com/ypapax/status_check/fake_service"
)

func main() {
	if err := logrus_conf.Files("fake-services", logrus.TraceLevel); err != nil {
		panic(err)
	}
	var confPath string
	flag.StringVar(&confPath, "conf", "fake-services.test.conf.yaml", "path to a config file")
	flag.Parse()
	f, err := os.Open(confPath)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	c, err := fake_config.Parse(f)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	if err := fake_service.Serve(c); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
