package psql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ypapax/jsn"

	"github.com/ypapax/status_check/status"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type statusRepository struct {
	db *sql.DB
}

func NewPostgresStatusRepository(db *sql.DB) status.Repository {
	return &statusRepository{
		db,
	}
}

func (r *statusRepository) Create(statuses []status.Status) error {
	const argsAmountPerRecord = 4
	var argsNumbers = make([]string, len(statuses))
	var args = make([]interface{}, len(statuses)*argsAmountPerRecord)
	var argCounter = 1
	for i, s := range statuses {
		argsNumbers[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)", argCounter, argCounter+1, argCounter+2, argCounter+3)
		args[argCounter-1] = s.Available
		args[argCounter] = s.Created
		args[argCounter+1] = int(s.ResponseTime.Seconds() * 1000)
		args[argCounter+2] = s.ServiceID
		argCounter += argsAmountPerRecord
	}
	query := fmt.Sprintf(`INSERT INTO status(available, created, response_ms, service_id) VALUES %s`, strings.Join(argsNumbers, ", "))
	result, err := r.db.Exec(query, args...)
	if err != nil {
		err := fmt.Errorf("error: %+v for query %+v, argNumbers: %+v", err, query, jsn.B(argsNumbers))
		logrus.Error(err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Printf("inserted %+v rows into status table, query: %+v, args: %+v", rows, query, jsn.B(args))
	return nil
}
func (r *statusRepository) AvailableServices(from, to time.Time) (int, error) {
	notAvail, err := r.NotAvailableServices(from, to)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	all, err := r.allServicesWithAnyStatus(from, to)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return all - notAvail, err
}

func (r *statusRepository) allServicesWithAnyStatus(from, to time.Time) (int, error) {
	query := `SELECT COUNT(DISTINCT id) FROM status WHERE created >= $1 AND created <= $2`
	var count = -1
	scan := func(rows *sql.Rows) error {
		if err := rows.Scan(&count); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	}
	if err := scanAll(r.db, scan, query, from, to); err != nil {
		logrus.Error(err)
		return 0, err
	}
	logrus.Tracef("count all services: %+v for query %+v with args %+v", count, query, from, to)
	return count, nil
}

func (r *statusRepository) NotAvailableServices(from, to time.Time) (int, error) {
	query := `SELECT COUNT(DISTINCT id) FROM status WHERE created >= $1 AND created <= $2 AND available = false`
	var count = -1
	scan := func(rows *sql.Rows) error {
		if err := rows.Scan(&count); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	}
	if err := scanAll(r.db, scan, query, from, to); err != nil {
		logrus.Error(err)
		return 0, err
	}
	logrus.Tracef("count not avail: %+v for query %+v with args %+v", count, query, from, to)
	return count, nil
}

func (r *statusRepository) ServicesWithResponseFasterThan(dur time.Duration, from, to time.Time) (int, error) {
	slow, err := r.ServicesWithResponseSlowerThan(dur, from, to)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	all, err := r.allServicesWithAnyStatus(from, to)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return all - slow, nil
}

func (r *statusRepository) ServicesWithResponseSlowerThan(dur time.Duration, from, to time.Time) (int, error) {
	query := `SELECT COUNT(DISTINCT id) FROM status WHERE created >= $1 AND created <= $2 AND available = true AND response_ms > $3`
	var count = -1
	scan := func(rows *sql.Rows) error {
		if err := rows.Scan(&count); err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	}
	if err := scanAll(r.db, scan, query, from, to, int(dur.Seconds()*1000)); err != nil {
		logrus.Error(err)
		return 0, err
	}
	return count, nil
}
