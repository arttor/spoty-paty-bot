package command

import (
	"fmt"
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type djLeftChat struct {
	next     Handler
	stateSvc *state.Service
	bot      *bot.BotAPI
}

func (s *djLeftChat) accepts(update bot.Update) bool {
	return update.Message != nil && update.Message.LeftChatMember != nil
}

func (s *djLeftChat) Handle(update bot.Update) () {
	err := s.stateSvc.Logout(update.Message.Chat, update.Message.LeftChatMember)
	if err == nil {
		_, err = s.bot.Send(bot.NewMessage(update.Message.Chat.ID, fmt.Sprintf(res.TxtLogoutSuccessPattern, update.Message.LeftChatMember.String())))
		if err != nil {
			logrus.WithError(err).Error("Unable to send logout response")
		}
	}
}
