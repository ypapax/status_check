package main

import (
	"flag"
	"github.com/ypapax/status_check"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
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
	if err := status_check.ServeAPI(*c); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
