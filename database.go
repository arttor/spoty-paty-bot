package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

var schema = `
CREATE TABLE IF NOT EXISTS chat (
    id bigint PRIMARY KEY,
    max_songs integer,
    songs integer,
    songs_user_id integer,
    spotify_token text
)
`

func setupDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.WithError(err).Fatal("Unable to connect DB")
	}
	db.MustExec(schema)
	return db
}

