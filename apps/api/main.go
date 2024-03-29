package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/logrus_conf"
	"github.com/ypapax/status_check/api"
	"github.com/ypapax/status_check/config"
)

func main() {
	if err := logrus_conf.Files("api", logrus.TraceLevel); err != nil {
		panic(err)
	}
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
	if err := api.Serve(*c); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
