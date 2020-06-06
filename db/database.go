package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

func New() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.WithError(err).Error("Unable to connect DB")
		return nil, err
	}
	_, err = db.Exec(createSearchTable)
	if err != nil {
		logrus.WithError(err).Error("Create schema")
		return nil, err
	}
	return db, nil
}
