package test

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

var serviceAddr string

const reqTimeout = 3 * time.Second

func TestMain(m *testing.M) {
	logrus.SetReportCaller(true)
	flag.StringVar(&serviceAddr, "service-addr", "http://localhost:3000", "address of status_check web service")
	flag.Parse()
	os.Exit(m.Run())
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
