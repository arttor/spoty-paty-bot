package service

import (
	"fmt"
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type loginFinish struct {
	stateSvc *state.Service
	bot      *bot.BotAPI
	next     Handler
}

func (s *loginFinish) Handle(update bot.Update) () {
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
func (s *loginFinish) accepts(update bot.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == res.CmdLoginFinish
}

func (s *loginFinish) handle(update bot.Update) {
	loginCode := update.Message.CommandArguments()
	err := s.stateSvc.FinishLogin(update, loginCode)
	var msg bot.MessageConfig
	if err != nil {
		msg = bot.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = bot.NewMessage(update.Message.Chat.ID, fmt.Sprintf(res.TxtFinishLoginSuccessPattern, update.Message.From.UserName))
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		logrus.WithError(err).Error("Unable to send already logged in response")
	}
}
