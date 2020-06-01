package state

import (
	"encoding/json"
	"errors"
	app "github.com/arttor/spoty-paty-bot/config"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
	"sync"
)

type Service struct {
	db     *sqlx.DB
	config app.Config
	m      sync.RWMutex
	mem    map[int64]Chat
}

func New(config app.Config, db *sqlx.DB) *Service {
	return &Service{config: config, db: db, mem: make(map[int64]Chat)}
}

func (s *Service) SaveToken(chatID int64, token *oauth2.Token) error {
	s.m.Lock()
	defer s.m.Unlock()
	chat, ok := s.mem[chatID]
	if !ok {
		return errors.New("no such chat")
	}
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}
	chat.SpotifyToken = string(tokenBytes)
	s.mem[chatID] = chat
	return nil
}

func (s *Service) AddIfExists(chat Chat) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.mem[chat.Id]
	if !ok {
		s.mem[chat.Id] = chat
	}
	return nil
}

func (s *Service) Get(chatID int64) (Chat,error) {
	s.m.RLock()
	defer s.m.RUnlock()
	chat, ok := s.mem[chatID]
	if !ok {
		return Chat{},errors.New("not found")
	}
	return chat,nil
}
