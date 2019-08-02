package statuses

import (
	"github.com/ypapax/status_check/status"
	"sync"
)

type Statuses struct {
	mtx      sync.Mutex
	statuses []status.Status
}

func (sts *Statuses) Add(st status.Status) {
	sts.mtx.Lock()
	defer sts.mtx.Unlock()
	sts.statuses = append(sts.statuses, st)
}

func (sts *Statuses) GetAll() []status.Status {
	sts.mtx.Lock()
	defer sts.mtx.Unlock()
	statuses := sts.statuses
	sts.statuses = nil
	return statuses
}
