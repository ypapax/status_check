package pipes

import (
	"context"
	"github.com/ypapax/status_check/status"
)

type StatusPipe interface {
	Publish(parent context.Context, status status.Status) error
	Listen(ctx context.Context, statusChan chan<- status.Status, errs chan<- error)
}
