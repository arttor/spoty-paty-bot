package service

import (
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type login struct {
	next       Handler
	stateSvc   *state.Service
	spotifySvc *spotify.Service
	bot        *bot.BotAPI
}

func (s *login) Handle(update tgbotapi.Update) () {
	if s.Accepts(update) {
		s.handle(update)
		return
	}
	if s.next != nil {
		s.next.Handle(update)
		return
	}
	logrus.Info("No handler for given update")
}
func (s *login) Accepts(update tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == "login"
}

func (s *login) handle(update tgbotapi.Update) {
	_, err := s.stateSvc.Get(update.Message.Chat.ID)
	if err != nil {
		s.login(update)
	} else {
		s.alreadyLoggedIn(update)
	}
}

func (s *login) login(update tgbotapi.Update) {
	url := s.spotifySvc.GetAuthURL(update.Message.Chat.ID)
	response := tgbotapi.NewMessage(update.Message.Chat.ID, url)
	_, err := s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send log in request")
	}
}

func (s *login) alreadyLoggedIn(update tgbotapi.Update) {
	response := tgbotapi.NewMessage(update.Message.Chat.ID, "Already logged in")
	_, err := s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send already logged in response")
	}
}
