package command

import (
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
)

type addSong struct {
	stateSvc *state.Service
	bot      *bot.BotAPI
}

func (s *addSong) accepts(update bot.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == res.CmdAddSong
}

func (s *addSong) Handle(update bot.Update) () {
	songID := update.Message.CommandArguments()
	if songID == "" {
		return
	}
	err := s.stateSvc.QueueSong(update.Message.From, update.Message.Chat, spotify.ID(songID))
	var msg bot.MessageConfig
	if err != nil {
		msg = bot.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = bot.NewMessage(update.Message.Chat.ID, res.TxtAddSongSuccess)
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		logrus.WithError(err).Error("Unable to send logout response")
	}
}
