package command

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
}

func (s *loginFinish) accepts(update bot.Update) bool {
	return update.Message!=nil && update.Message.IsCommand() && update.Message.Command() == res.CmdLoginFinish
}

func (s *loginFinish) Handle(update bot.Update) () {
	_, _ = s.bot.DeleteMessage(bot.DeleteMessageConfig{
	ChatID:    update.Message.Chat.ID,
	MessageID: update.Message.MessageID,
})
	loginCode := update.Message.CommandArguments()
	err := s.stateSvc.FinishLogin(update, loginCode)
	var msg bot.MessageConfig
	if err != nil {
		msg = bot.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = bot.NewMessage(update.Message.Chat.ID, fmt.Sprintf(res.TxtFinishLoginSuccessPattern, update.Message.From.String()))
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		logrus.WithError(err).Error("Unable to send finish login response")
	}
}
