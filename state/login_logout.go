package state

import (
	"errors"
	"fmt"
	app "github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/res"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify"
	"math/rand"
	"sync"
)

const loginCodeLength = 8

type Service struct {
	db     *sqlx.DB
	config app.Config
	m      sync.RWMutex
	mem    map[int64]Chat
}

func New(config app.Config, db *sqlx.DB) *Service {
	return &Service{config: config, db: db, mem: make(map[int64]Chat)}
}

func (s *Service) SaveClient(chatID int64, client *spotify.Client) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()
	chat, ok := s.mem[chatID]
	if !ok {
		return "", errors.New(res.TxtWebLoginNoSuchChat)
	}
	if chat.DjID != 0 {
		return "", errors.New(res.TxtWebLoginAlready)
	}
	chat.DjClient = client
	if chat.LoginCandidates == nil {
		chat.LoginCandidates = make(map[string]*spotify.Client)
	}
	loginCode := randSeq(loginCodeLength)
	chat.LoginCandidates[loginCode] = client
	s.mem[chatID] = chat
	return loginCode, nil
}

func (s *Service) AddIfNotExists(chat Chat) {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.mem[chat.Id]
	if !ok {
		s.mem[chat.Id] = chat
	}
}

func (s *Service) Get(chatID int64) (Chat, bool) {
	s.m.RLock()
	defer s.m.RUnlock()
	chat, ok := s.mem[chatID]
	return chat, ok
}

func (s *Service) Logout(fromChat *bot.Chat, userLeft *bot.User) error {
	s.m.Lock()
	defer s.m.Unlock()
	chatID := fromChat.ID
	chat, ok := s.mem[chatID]
	if !ok || chat.DjID == 0 {
		return errors.New(res.TxtLogoutErrNotLogin)
	}
	if chat.DjID != userLeft.ID {
		return errors.New(fmt.Sprintf(res.TxtLogoutErrAnotherUserPattern, chat.DjName))
	}
	chat.DjClient = nil
	chat.LoginCandidates = nil
	chat.DjID = 0
	chat.DjName = ""
	chat.Queue = nil
	s.mem[chatID] = chat
	return nil
}

func (s *Service) FinishLogin(update bot.Update, code string) error {
	s.m.Lock()
	defer s.m.Unlock()
	chatID := update.Message.Chat.ID
	chat, ok := s.mem[chatID]
	if !ok {
		return errors.New(res.TxtLogoutErrNotLogin)
	}
	if chat.DjID != 0 {
		return errors.New(fmt.Sprintf(res.TxtLoginAlreadyPattern, chat.DjName))
	}
	if chat.LoginCandidates == nil || chat.LoginCandidates[code] == nil {
		return errors.New(res.TxtLogoutErrNotLogin)
	}
	chat.DjClient = chat.LoginCandidates[code]
	chat.DjID = update.Message.From.ID
	chat.DjName = update.Message.From.String()
	chat.LoginCandidates = nil
	s.mem[chatID] = chat
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
