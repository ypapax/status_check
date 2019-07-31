package status

import "time"

type Service interface {
	CreateStatus(status *Status, serviceID int) error
	AvailableServices(from, to time.Time) (int, error)
	NotAvailableServices(from, to time.Time) (int, error)
	ServicesWithResponseFasterThan(dur time.Duration, from, to time.Time) (int, error)
	ServicesWithResponseSlowerThan(dur time.Duration, from, to time.Time) (int, error)
}

type statusService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &statusService{
		repo: repo,
	}
}

func (s *statusService) CreateStatus(status *Status, serviceID int) error {
	status.Created = time.Now()
	status.ServiceID = serviceID
	return s.repo.Create(status)
}

func (s *statusService) AvailableServices(from, to time.Time) (int, error) {
	return s.repo.AvailableServices(from, to)
}

func (s *statusService) NotAvailableServices(from, to time.Time) (int, error) {
	return s.repo.AvailableServices(from, to)
}

func (s *statusService) ServicesWithResponseFasterThan(dur time.Duration, from, to time.Time) (int, error) {
	return s.repo.ServicesWithResponseFasterThan(dur, from, to)
}

func (s *statusService) ServicesWithResponseSlowerThan(dur time.Duration, from, to time.Time) (int, error) {
	return s.repo.ServicesWithResponseSlowerThan(dur, from, to)
}
