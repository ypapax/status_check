package job

import (
	web_service "github.com/ypapax/status_check/web-service"
	"time"
)

type Job struct {
	Service web_service.WebService
	LastCheckedTime time.Time
}
