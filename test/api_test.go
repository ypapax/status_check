package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAvailable(t *testing.T) {
	as := assert.New(t)
	from := time.Now().Add(-2 * time.Minute)
	to := time.Now().Add(-time.Minute)
	status, b, err := getPath(fmt.Sprintf("/services-count/available/%d/%d", from.Unix(), to.Unix()))
	if !as.NoError(err) {
		return
	}
	if !as.Equal(status, http.StatusOK) {
		return
	}
	c, err := getCount(b)
	if !as.NoError(err) {
		return
	}
	if as.True(c > 0, "count is zero") {
		return
	}
}
