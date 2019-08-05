package queue_test

import (
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ypapax/status_check/queue"
)

func TestNext(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
	as := assert.New(t)
	q := queue.New([]int{1, 2})
	as.Equal(1, q.Next())
	as.Equal(2, q.Next())
	as.Equal(1, q.Next())
}
