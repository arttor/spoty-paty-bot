package service

import (
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type logout struct {
	next       Handler
	stateSvc   *state.Service
	bot        *bot.BotAPI
}

func (s *logout) Handle(update bot.Update) () {
	if s.accepts(update) {
		s.handle(update)
		return
	}
	if s.next != nil {
		s.next.Handle(update)
		return
	}
	logrus.Info("No handler for given update")
}
func (s *logout) accepts(update bot.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == "logout"
}

func (s *logout) handle(update bot.Update) {
	s.stateSvc.RemoveClient(update.Message.Chat.ID)
	response := bot.NewMessage(update.Message.Chat.ID, "Logged out")
	_, err := s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send already logged in response")
	}
}
