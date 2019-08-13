package status_check

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ypapax/status_check/database/psql"
	"github.com/ypapax/status_check/status"
	web_service "github.com/ypapax/status_check/web-service"
)

func Services(dbType, connString string) (web_service.Service, status.Service, error) {
	var webServicesService web_service.Service
	var statusService status.Service

	switch dbType {
	case "psql":
		db, err := psql.GetConnection(connString, 10*time.Second)
		if err != nil {
			logrus.Error(err)
			return webServicesService, statusService, err
		}
		serviceRepo := psql.NewPostgresServiceRepository(db)
		statusRepo := psql.NewPostgresStatusRepository(db)
		webServicesService = web_service.NewService(serviceRepo)
		statusService = status.NewService(statusRepo)
	default:
		err := fmt.Errorf("db type '%+v' is not supported", connString)
		return webServicesService, statusService, err
	}
	return webServicesService, statusService, nil
}
