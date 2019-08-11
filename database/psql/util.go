package psql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func GetConnection(database string, timeout time.Duration) (*sql.DB, error) {
	logrus.Println("Connecting to PostgreSQL DB: ", database)
	db, err := sql.Open("postgres", database)
	if err != nil {
		err := fmt.Errorf("couldn't connect to db: %+v, err: %+v", database, err)
		logrus.Error(err)
		return nil, err
	}
	done := make(chan bool)
	go func() {
		for {
			if err := db.Ping(); err != nil {
				logrus.Println("waiting for database ", database)
				time.Sleep(2 * time.Second)
			}
			done <- true
			break
		}
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		err := fmt.Errorf("couldn't ping db %+v after timeout %s", database, timeout)
		return nil, err

	}
	return db, nil
}

func scanAll(db *sql.DB, scan func(rows *sql.Rows) error, query string, args ...interface{}) error {
	if db == nil {
		err := fmt.Errorf("db is nil")
		logrus.Error(err)
		return err
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := scan(rows); err != nil {
			logrus.Error(err)
			return err
		}

	}
	return nil
}
