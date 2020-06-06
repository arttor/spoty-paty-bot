package command

import (
	"github.com/arttor/spoty-paty-bot/search"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type searchLogout struct {
	searchSvc search.Service
	bot      *bot.BotAPI
	command  string
}

func (s *searchLogout) Handle(update bot.Update) () {
	s.searchSvc.Logout()
	_, err := s.bot.Send( bot.NewMessage(update.Message.Chat.ID, "Logout success."))
	if err != nil {
		logrus.WithError(err).Error("Unable to send logout response")
	}
}
func (s *searchLogout) accepts(update bot.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == s.command
}
