package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/sirupsen/logrus"

	"github.com/ypapax/status_check"
	"github.com/ypapax/status_check/fake_service"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(time.Hour)

	limitMS := 1000

	type pathAndExpected struct {
		path          string
		expectedCount int
	}

	type testCase struct {
		name            string
		paths           []pathAndExpected
		statusCheckConf status_check.Config
		fakeServiceConf fake_service.Config
		workTime        time.Duration
	}

	cases := []testCase{
		{
			name:     "simple",
			workTime: 20 * time.Second,
			paths: []pathAndExpected{
				{path: fmt.Sprintf("/services-count/available/%d/%d", from.Unix(), to.Unix()), expectedCount: 1},
				{path: fmt.Sprintf("/services-count/not-available/%d/%d", from.Unix(), to.Unix()), expectedCount: 1},
				{path: fmt.Sprintf("/services-count/faster/%d/%d/%d", limitMS, from.Unix(), to.Unix()), expectedCount: 2},
				{path: fmt.Sprintf("/services-count/slower/%d/%d/%d", limitMS, from.Unix(), to.Unix()), expectedCount: 0},
			},
			statusCheckConf: status_check.Config{
				Bind:             "0.0.0.0:3001",
				CheckPeriod:      5 * time.Second,
				DbType:           "psql",
				ConnectionString: "postgresql://postgres@postgres/status_check?sslmode=disable",
				Workers:          100,
				Schemas:          []string{"https", "http"},
			},
			fakeServiceConf: fake_service.Config{
				Ports: []fake_service.Port{
					{From: 2001, To: 2001, StatusCodes: []int{200}},
					{From: 3001, To: 3001, StatusCodes: []int{502}},
				},
			},
		},
		{
			name:     "diff_status",
			workTime: 20 * time.Second,
			paths: []pathAndExpected{
				{path: fmt.Sprintf("/services-count/available/%d/%d", from.Unix(), to.Unix()), expectedCount: 1},
				{path: fmt.Sprintf("/services-count/not-available/%d/%d", from.Unix(), to.Unix()), expectedCount: 1},
				{path: fmt.Sprintf("/services-count/faster/%d/%d/%d", limitMS, from.Unix(), to.Unix()), expectedCount: 2},
				{path: fmt.Sprintf("/services-count/slower/%d/%d/%d", limitMS, from.Unix(), to.Unix()), expectedCount: 0},
			},
			statusCheckConf: status_check.Config{
				Bind:             "0.0.0.0:3001",
				CheckPeriod:      5 * time.Second,
				DbType:           "psql",
				ConnectionString: "postgresql://postgres@postgres/status_check?sslmode=disable",
				Workers:          100,
				Schemas:          []string{"https", "http"},
			},
			fakeServiceConf: fake_service.Config{
				Ports: []fake_service.Port{
					{From: 2001, To: 2001, StatusCodes: []int{200}},
					{From: 3001, To: 3001, StatusCodes: []int{200, 502}},
				},
			},
		},
		{
			name:     "big",
			workTime: 60 * time.Second,
			paths: []pathAndExpected{
				{path: fmt.Sprintf("/services-count/available/%d/%d", from.Unix(), to.Unix()), expectedCount: 1021},
				{path: fmt.Sprintf("/services-count/not-available/%d/%d", from.Unix(), to.Unix()), expectedCount: 13},
				{path: fmt.Sprintf("/services-count/faster/%d/%d/%d", limitMS, from.Unix(), to.Unix()), expectedCount: 1014},
				{path: fmt.Sprintf("/services-count/slower/%d/%d/%d", limitMS, from.Unix(), to.Unix()), expectedCount: 20},
			},
			statusCheckConf: status_check.Config{
				Bind:             "0.0.0.0:3001",
				CheckPeriod:      20 * time.Second,
				DbType:           "psql",
				ConnectionString: "postgresql://postgres@postgres/status_check?sslmode=disable",
				Workers:          200,
				Schemas:          []string{"https", "http"},
			},
			fakeServiceConf: fake_service.Config{
				Ports: []fake_service.Port{
					{From: 2001, To: 3000, StatusCodes: []int{200}, DelayMS: []int{500}},
					{From: 8001, To: 8010, StatusCodes: []int{200}, DelayMS: []int{1100}},
					{From: 9001, To: 9010, StatusCodes: []int{200}, DelayMS: []int{1100, 500}},
					{From: 7001, To: 7010, StatusCodes: []int{200, 502}, DelayMS: []int{200, 502}},
					{From: 3001, To: 3001, StatusCodes: []int{501}},
					{From: 4001, To: 4001, StatusCodes: []int{502}},
					{From: 5001, To: 5001, StatusCodes: []int{503}},
					{From: 6001, To: 6001, StatusCodes: []int{504}},
				},
			},
		},
	}

	for _, c := range cases {
		func() {
			c.statusCheckConf.Addresses = allFakeServicesAddr(fakeServicesContainerName, c.fakeServiceConf)
			b, err := yaml.Marshal(c.fakeServiceConf)
			if err != nil {
				logrus.Error(err)
			}
			f := `../fake-services.test.conf.yaml`
			if err := ioutil.WriteFile(f, b, 0777); err != nil {
				logrus.Error(err)
			}

			stopContainers, err := launchContainers(c.statusCheckConf)
			if err != nil {
				logrus.Error(err)
			}
			defer func() {
				if err := stopContainers(); err != nil {
					logrus.Error(err)
				}
			}()
			logrus.Infof("waiting %s before running tests", c.workTime)
			time.Sleep(c.workTime)
			for _, p := range c.paths {
				t.Run(c.name+p.path, func(t *testing.T) {
					as := assert.New(t)
					status, b, err := getPath(p.path)
					t.Log("resp: ", string(b))
					if !as.NoError(err) {
						return
					}
					if !as.Equal(http.StatusOK, status) {
						return
					}
					count, err := getCount(b)
					if !as.NoError(err) {
						return
					}
					if as.Equal(p.expectedCount, count) {
						return
					}
				})
			}
		}()
	}
}
