package command

import (
	"github.com/arttor/spoty-paty-bot/inlinesearch"
	"github.com/arttor/spoty-paty-bot/res"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type logout struct {
	stateSvc inlinesearch.Service
	bot      *bot.BotAPI
}

func (s *logout) Handle(update bot.Update) () {
	s.stateSvc.Logout()
	_, err := s.bot.Send( bot.NewMessage(update.Message.Chat.ID, "Logout success."))
	if err != nil {
		logrus.WithError(err).Error("Unable to send logout response")
	}
}
func (s *logout) accepts(update bot.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == res.CmdLogout
}
