package search

import (
	"encoding/json"
	"github.com/arttor/spoty-paty-bot/db"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type Service interface {
	SetClient(spotify.Client)
	Close()
	RestoreClient(getClient GetClient)
	Logout()
	GetClient()*spotify.Client
}

type GetClient func(token *oauth2.Token) spotify.Client

type service struct {
	client    *spotify.Client
	db        *sqlx.DB
	getClient GetClient
}

func NewService(database *sqlx.DB) Service {
	return &service{db: database}
}

func (s *service) SetClient(client spotify.Client) {
	s.client = &client
}

func (s *service) Close() {
	if s.client == nil || s.db == nil {
		return
	}
	token, err := s.client.Token()
	if err != nil {
		logrus.Error("Persist state: get token from client error")
		return
	}
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		logrus.Error("Persist state: unable to marshal token")
		return
	}
	searchToken := db.Search{}
	err = s.db.Get(&searchToken, "SELECT * FROM search WHERE id=$1", true)
	isInsert := err != nil
	tx, err := s.db.Beginx()
	if err != nil {
		logrus.Error(err)
		return
	}
	if isInsert {
		_, err = tx.Exec("INSERT INTO search (id, token) VALUES ($1, $2)", true, string(tokenBytes))
	} else {
		_, err = tx.Exec("UPDATE search SET token = $1 WHERE id = $2", string(tokenBytes), true)
	}
	if err != nil {
		logrus.Error(err)
		_ = tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return
	}

}

func (s *service) RestoreClient(getClient GetClient) {
	searchToken := db.Search{}
	err := s.db.Get(&searchToken, "SELECT * FROM search WHERE id=$1", true)
	if err != nil {
		logrus.WithError(err).Info("No search config in db to restore")
		return
	}
	token := oauth2.Token{}
	err = json.Unmarshal([]byte(searchToken.Token), &token)
	if err != nil {
		logrus.WithError(err).Infof("unable to unmarshal token %v", searchToken.Token)
		return
	}
	res := getClient(&token)
	s.client = &res
}

func (s *service) Logout() {
	tx, err := s.db.Beginx()
	if err != nil {
		logrus.Error(err)
		return
	}
	_, err = tx.Exec("DELETE FROM search WHERE id = $1", true)
	if err != nil {
		logrus.Error(err)
		_ = tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return
	}
	s.client = nil
}

func (s *service) GetClient() *spotify.Client {
	return s.client
}
