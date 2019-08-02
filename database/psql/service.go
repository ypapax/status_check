package psql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	web_service "github.com/ypapax/status_check/web-service"
)

type serviceRepository struct {
	db *sql.DB
}

func NewPostgresServiceRepository(db *sql.DB) web_service.Repo {
	return &serviceRepository{
		db,
	}
}

func (r *serviceRepository) Create(services []web_service.WebService) error {
	var args []interface{}
	var argNumbers = make([]string, len(services))
	var argCounter = 1
	for i, s := range services {
		args = append(args, s.Address)
		args = append(args, s.Created)
		argNumbers[i] = fmt.Sprintf("($%d, $%d)", argCounter, argCounter+1)
		argCounter += 2
	}
	query := fmt.Sprintf(`INSERT INTO service(address, created)
		VALUES %s ON CONFLICT (address) DO NOTHING`, strings.Join(argNumbers, ", "))
	if _, err := r.db.Exec(query, args...); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (r *serviceRepository) FindAll() ([]web_service.WebService, error) {
	var ss []web_service.WebService
	query := `SELECT id, address, created FROM service`
	scan := func(rows *sql.Rows) error {
		var s web_service.WebService
		if err := rows.Scan(&s.ID, &s.Address, &s.Created); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	}
	if err := scanAll(r.db, scan, query); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return ss, nil
}
