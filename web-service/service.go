package web_service

import "time"

type Service interface {
	CreateWebService(service *WebService) error
}

type webServiceService struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &webServiceService{
		repo: repo,
	}
}

func (s *webServiceService) CreateWebService(service *WebService) error {
	service.Created = time.Now()
	return s.repo.Create(service)
}
