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
	paths := []string{
		fmt.Sprintf("/services-count/available/%d/%d", from.Unix(), to.Unix()),
		fmt.Sprintf("/services-count/not-available/%d/%d", from.Unix(), to.Unix()),
		fmt.Sprintf("/services-count/faster/%d/%d/%d", limitMS, from.Unix(), to.Unix()),
		fmt.Sprintf("/services-count/slower/%d/%d/%d", limitMS, from.Unix(), to.Unix()),
	}
	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			as := assert.New(t)
			status, b, err := getPath(p)
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
