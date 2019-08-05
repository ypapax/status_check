package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/ypapax/status_check"
	"gopkg.in/yaml.v2"

	"github.com/ypapax/status_check/fake_service"

	"github.com/sirupsen/logrus"
)

var serviceAddr string
var fakeServicesConfFile string
var fakeServicesContainerName string
var dockerComposeConfigFile string
var waitBeforeRunningTestsSeconds int

const reqTimeout = 3 * time.Second

func TestMain(m *testing.M) {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
	flag.StringVar(&serviceAddr, "service-addr", "http://localhost:3001", "address of status_check web service")
	flag.StringVar(&dockerComposeConfigFile, "docker-compose", "../docker-compose-test.yml", "docker compose config file")
	flag.StringVar(&fakeServicesContainerName, "fake-container", "fake-services", "fake services container name")
	flag.IntVar(&waitBeforeRunningTestsSeconds, "wait-secs", 60, "amount of seconds to wait when status_check service collects enough stats before running tests")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}

func launchContainers(statusCheckConf status_check.Config) (func() error, error) {
	b, err := yaml.Marshal(statusCheckConf)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	statusCheckConfFile := "./status_check.test.conf.yaml"
	logrus.Tracef("writing file %+v", statusCheckConfFile)
	if err := ioutil.WriteFile(statusCheckConfFile, b, 0777); err != nil {
		logrus.Error(err)
		return nil, err
	}

	buildCompose := exec.Command(`docker-compose`, "-f", dockerComposeConfigFile, "build")
	logrus.Tracef("running: %+v", strings.Join(buildCompose.Args, " "))
	buildCompose.Stderr = os.Stderr
	buildCompose.Stdout = os.Stdout
	if err := buildCompose.Run(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	runCompose := exec.Command(`docker-compose`, "-f", dockerComposeConfigFile, "up", "-d", "--force-recreate")
	logrus.Tracef("running: %+v", strings.Join(runCompose.Args, " "))
	runCompose.Stderr = os.Stderr
	runCompose.Stdout = os.Stdout
	if err := runCompose.Run(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	for {
		_, err := http.Get(serviceAddr)
		if err != nil {
			w := 10 * time.Second
			logrus.Infof("waiting for %+v, sleeping for %s", serviceAddr, w)
			time.Sleep(w)
			continue
		}
		break
	}
	return func() error {
		downCompose := exec.Command(`docker-compose`, "-f", dockerComposeConfigFile, "down")
		logrus.Tracef("running: %+v", strings.Join(downCompose.Args, " "))
		downCompose.Stderr = os.Stderr
		downCompose.Stdout = os.Stdout
		if err := downCompose.Run(); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	}, nil
}

func getPath(path string) (int, []byte, error) {
	u := serviceAddr + path
	logrus.Println("requesting ", u)
	var netClient = &http.Client{
		Timeout: reqTimeout,
	}
	response, err := netClient.Get(u)
	if err != nil {
		logrus.Error(err)
		return 0, nil, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Error(err)
		return response.StatusCode, nil, err
	}
	return response.StatusCode, b, nil
}

type countResp struct {
	Count int
}

func getCount(resp []byte) (int, error) {
	var c countResp
	if err := json.Unmarshal(resp, &c); err != nil {
		logrus.Error(err)
		return 0, err
	}
	return c.Count, nil
}

func allFakeServicesAddr(fakeServicesContainerName string, conf fake_service.Config) []string {
	ports := fake_service.AllPorts(&conf)
	var fakeServicesAddr []string
	for _, p := range ports {
		fakeServicesAddr = append(fakeServicesAddr, fmt.Sprintf("%s:%d", fakeServicesContainerName, p))
	}
	return fakeServicesAddr
}
