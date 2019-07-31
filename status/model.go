package status

import "time"

type Status struct {
	Available bool
	ResponseTime time.Duration
	Created time.Time
	ServiceID int
}
