package command

import (
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type queue struct {
	stateSvc *state.Service
	bot      *bot.BotAPI
}

func (s *queue) accepts(update bot.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == res.CmdQueue
}

func (s *queue) Handle(update bot.Update) () {
	//TODO: implement
	_, _ = s.stateSvc.GetQueue(update.Message.Chat)
}
