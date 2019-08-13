package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check/checker"
	"github.com/ypapax/status_check/config"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)
	var confPath string
	flag.StringVar(&confPath, "conf", "conf.yaml", "path to config file")
	flag.Parse()
	f, err := os.Open(confPath)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	c, err := config.Parse(f)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	if err := checker.Check(*c); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

}
