package service

import (
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type logout struct {
	next     Handler
	stateSvc *state.Service
	bot      *bot.BotAPI
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
	return update.Message!=nil && update.Message.IsCommand() && update.Message.Command() == res.CmdLogout
}

func (s *logout) handle(update bot.Update) {
	err := s.stateSvc.Logout(update)
	var msg bot.MessageConfig
	if err != nil {
		msg = bot.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = bot.NewMessage(update.Message.Chat.ID, res.TxtLogoutSuccess)
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		logrus.WithError(err).Error("Unable to send logout response")
	}
}
