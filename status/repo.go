package status

import "time"

type Repository interface {
	Create(status []Status) error
	AvailableServices(from, to time.Time) (int, error)
	NotAvailableServices(from, to time.Time) (int, error)
	ServicesWithResponseFasterThan(dur time.Duration, from, to time.Time) (int, error)
	ServicesWithResponseSlowerThan(dur time.Duration, from, to time.Time) (int, error)
}
