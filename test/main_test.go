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

	"github.com/ypapax/status_check/fake_config"

	"github.com/ypapax/status_check/config"

	"gopkg.in/yaml.v2"

	"github.com/ypapax/status_check/fake_service"

	"github.com/sirupsen/logrus"
)

var serviceAddr string
var fakeServicesConfFile string
var fakeServicesContainerName string
var dockerComposeConfigFile string

const reqTimeout = 3 * time.Second

func TestMain(m *testing.M) {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
	flag.StringVar(&serviceAddr, "service-addr", "http://localhost:3001", "address of status_check web service")
	flag.StringVar(&dockerComposeConfigFile, "docker-compose", "../docker-compose-test.yml", "docker compose config file")
	flag.StringVar(&fakeServicesContainerName, "fake-container", "fake-services", "fake services container name")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}

func launchContainers(statusCheckConf config.Config) (func() error, error) {
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
	createTestNetwork := exec.Command(`docker`, `network`, `create`, `test-network`)
	logrus.Tracef("running: %+v", strings.Join(createTestNetwork.Args, " "))
	createTestNetwork.Stderr = os.Stderr
	createTestNetwork.Stdout = os.Stdout
	if err := createTestNetwork.Run(); err != nil {
		logrus.Warn(err)
	}
	buildCompose := exec.Command(`docker-compose`, "-f", dockerComposeConfigFile, "build")
	logrus.Tracef("running: %+v", strings.Join(buildCompose.Args, " "))
	buildCompose.Stderr = os.Stderr
	buildCompose.Stdout = os.Stdout
	if err := buildCompose.Run(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	upOutFileName := fmt.Sprintf("/tmp/%d-status-check-docker-compose-up.stderr.stdout", time.Now().Unix())
	f, err := os.Create(upOutFileName)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	ps := func() {
		dockerPs := exec.Command("docker", "ps")
		logrus.Tracef("running: %+v", strings.Join(dockerPs.Args, " "))
		dockerPs.Stderr = os.Stderr
		dockerPs.Stdout = os.Stdout
		if err := dockerPs.Run(); err != nil {
			logrus.Error(err)
		}
	}
	runCompose := exec.Command(`docker-compose`, "-f", dockerComposeConfigFile, "up", "--force-recreate")
	shutdown := func() error {
		if runCompose.Process != nil {
			if err := runCompose.Process.Kill(); err != nil {
				logrus.Error(err)
			}
		}
		defer f.Close()
		ps()
		downCompose := exec.Command(`docker-compose`, "-f", dockerComposeConfigFile, "down", "-v")
		logrus.Tracef("running: %+v", strings.Join(downCompose.Args, " "))
		downCompose.Stderr = os.Stderr
		downCompose.Stdout = os.Stdout
		if err := downCompose.Run(); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	}
	outFileContent := func() string {
		b, err := ioutil.ReadFile(upOutFileName)
		if err != nil {
			logrus.Error(err)
			return err.Error()
		}
		return fmt.Sprintf(`
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
content of the file %+v: %s\n
\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/
`, upOutFileName, string(b))
	}
	errChan := make(chan error)
	go func() {
		logrus.Tracef("running: %+v, stdout and stderr are written to %s", strings.Join(runCompose.Args, " "), upOutFileName)
		runCompose.Stderr = f
		runCompose.Stdout = f
		if err := runCompose.Run(); err != nil {
			err := fmt.Errorf("error: %+v, for command %s", err, strings.Join(runCompose.Args, " "))
			logrus.Warn(outFileContent())
			logrus.Error(err)
			//errChan <- err
			return
		}
	}()
	ready := make(chan bool, 2)
	go func() {
		for {
			_, err := http.Get(serviceAddr)
			if err != nil {
				/*runStatusCheckCompose := exec.Command(`docker-compose`, "-f", dockerComposeConfigFile, "up", "-d", "api")
				logrus.Tracef("running: %+v", strings.Join(runStatusCheckCompose.Args, " "))
				runStatusCheckCompose.Stderr = os.Stderr
				runStatusCheckCompose.Stdout = os.Stdout
				if err := runStatusCheckCompose.Run(); err != nil {
					logrus.Error(err)
					errChan <- err
					return
				}*/

				w := 5 * time.Second
				ps()
				logrus.Infof(`waiting for %+v, sleeping for %s, %+v`, serviceAddr, w, outFileContent())
				time.Sleep(w)
				continue
			}
			ready <- true
			break
		}
	}()
	dockerComposeUpTimeout := 25 * time.Second
	select {
	case <-ready:
		break
	case err := <-errChan:
		logrus.Error(err)
		return shutdown, err
	case <-time.After(dockerComposeUpTimeout):
		err := fmt.Errorf("timeout with launching containers := %s, %s", dockerComposeUpTimeout, outFileContent())
		logrus.Error(err)
		return shutdown, err
	}
	ps()
	return shutdown, nil
}

func getPath(path string, t *testing.T) (int, []byte, error) {
	u := serviceAddr + path
	t.Log("requesting ", u)
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

func allFakeServicesAddr(fakeServicesContainerName string, conf fake_config.Config) []string {
	ports := fake_service.AllPorts(&conf)
	var fakeServicesAddr []string
	for _, p := range ports {
		fakeServicesAddr = append(fakeServicesAddr, fmt.Sprintf("%s:%d", fakeServicesContainerName, p))
	}
	return fakeServicesAddr
}
