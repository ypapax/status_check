package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check"
	"os"
)

func main() {
	logrus.SetReportCaller(true)

	var confPath string
	flag.StringVar(&confPath, "conf", "conf.yaml", "path to config file")
	flag.Parse()
	f, err := os.Open(confPath)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	c, err := status_check.ParseConf(f)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	if err := status_check.Serve(*c); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
