package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	from := time.Now().Add(-2 * time.Minute)
	to := time.Now().Add(-time.Minute)

	limitMS := 1000

	type testCase struct {
		path          string
		expectedCount int
	}
	cases := []testCase{
		{path: fmt.Sprintf("/services-count/available/%d/%d", from.Unix(), to.Unix()), expectedCount: 1021},
		{path: fmt.Sprintf("/services-count/not-available/%d/%d", from.Unix(), to.Unix()), expectedCount: 13},
		{path: fmt.Sprintf("/services-count/faster/%d/%d/%d", limitMS, from.Unix(), to.Unix()), expectedCount: 1001},
		{path: fmt.Sprintf("/services-count/slower/%d/%d/%d", limitMS, from.Unix(), to.Unix()), expectedCount: 20},
	}
	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			as := assert.New(t)
			status, b, err := getPath(c.path)
			t.Log("resp: ", string(b))
			if !as.NoError(err) {
				return
			}
			if !as.Equal(http.StatusOK, status) {
				return
			}
			c, err := getCount(b)
			if !as.NoError(err) {
				return
			}
			t.Log("count: ", c)
			if as.True(c > 0, "count is zero") {
				return
			}
		})
	}
}
