package command

import (
	"fmt"
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type logout struct {
	stateSvc *state.Service
	bot      *bot.BotAPI
}

func (s *logout) Handle(update bot.Update) () {
	err := s.stateSvc.Logout(update.Message.Chat, update.Message.From)
	var msg bot.MessageConfig
	if err != nil {
		msg = bot.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = bot.NewMessage(update.Message.Chat.ID, fmt.Sprintf(res.TxtLogoutSuccessPattern,update.Message.From.String()))
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		logrus.WithError(err).Error("Unable to send logout response")
	}
}
func (s *logout) accepts(update bot.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == res.CmdLogout
}
