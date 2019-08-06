package web_service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check/status"
)

const reqTimeout = time.Second * 10

var notAvailableStatusCodes = map[int]struct{}{502: {}, 503: {}, 504: {}}

type Service interface {
	CreateWebServices(service []WebService) error
	CheckStatus(service *WebService, schemes []string) (*status.Status, error)
	FindAllWebServices() ([]WebService, error)
}

type webServiceService struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &webServiceService{
		repo: repo,
	}
}

func (s *webServiceService) CreateWebServices(services []WebService) error {
	now := time.Now()
	for _, s := range services {
		s.Created = now
	}
	return s.repo.Create(services)
}

func (s *webServiceService) FindAllWebServices() ([]WebService, error) {
	return s.repo.FindAll()
}

func (s *webServiceService) CheckStatus(service *WebService, schemes []string) (*status.Status, error) {
	logrus.Println("checking status for web service ", service, " with schemes ", schemes)
	if len(schemes) == 0 {
		err := fmt.Errorf("missing schemes")
		logrus.Error(err)
		return nil, err
	}
	var st *status.Status
	for _, sch := range schemes {
		code, dur, err := req(sch, service.Address)
		if err != nil {
			logrus.Error(err)
			continue
		}
		st = &status.Status{
			ResponseTime: *dur,
			ServiceID:    service.ID,
			Created:      time.Now(),
		}
		logrus.Println("status code ", *code, " for requesting ", sch, " ", service.Address)
		_, badStatus := notAvailableStatusCodes[*code]
		st.Available = !badStatus
		break
	}
	if st == nil {
		err := fmt.Errorf("couldn't get status for %+v", service)
		logrus.Error(err)
		return nil, err
	}
	return st, nil
}

func req(scheme, addr string) (*int, *time.Duration, error) {
	u := scheme + "://" + addr
	var netClient = &http.Client{
		Timeout: reqTimeout,
	}
	logrus.Println("requesting ", u, " with scheme ", scheme, " and addr ", addr)
	t1 := time.Now()
	response, err := netClient.Get(u)
	dur := time.Since(t1)
	if err != nil {
		logrus.Error(err)
		return nil, &dur, err
	}
	return &response.StatusCode, &dur, nil
}
