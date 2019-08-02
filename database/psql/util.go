package psql

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

func GetConnection(database string) (*sql.DB, error) {
	logrus.Println("Connecting to PostgreSQL DB")
	db, err := sql.Open("postgres", database)
	if err != nil {
		logrus.Error(err)
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