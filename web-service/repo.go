package web_service

type Repo interface {
	Create(service []WebService) error
	FindAll() ([]WebService, error)
}
