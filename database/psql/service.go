package psql

import (
	"database/sql"
	_ "github.com/lib/pq"
	web_service "github.com/ypapax/status_check/web-service"
)

type serviceRepository struct {
	db *sql.DB
}

func NewPostgresServiceRepository(db *sql.DB) web_service.Service{
	return &serviceRepository{
		db,
	}
}

func (r *serviceRepository) Create(service *web_service.WebService) error {
		if _, err := r.db.QueryRow(`INSERT INTO service(address, created)
		VALUES ($1) RETURNING id`,
			service.Address, service.Created).Scan(&service.ID); err != nil {
				return err
		}
	return nil
}
