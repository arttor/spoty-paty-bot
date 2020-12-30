package command

import (
	"fmt"
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"strings"
)

type addSong struct {
	stateSvc *state.Service
	bot      *bot.BotAPI
}

func (s *addSong) accepts(update bot.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == res.CmdAddSong
}

func (s *addSong) Handle(update bot.Update) () {
	_, err := s.bot.DeleteMessage(bot.DeleteMessageConfig{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.MessageID,
	})
	if err != nil {
		logrus.WithError(err).Error("Unable to delete message")
	}
	args := strings.Split(update.Message.CommandArguments(), "|")
	if len(args) < 1 {
		logrus.WithError(err).Error("Invalid addSong command arguments")
		return
	}
	err = s.stateSvc.AddSong(update.Message.From, update.Message.Chat, spotify.ID(args[0]))
	songName := ""
	if len(args) > 1 {
		songName = args[1]
	}
	var msg bot.MessageConfig
	if err != nil {
		msg = bot.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = bot.NewMessage(update.Message.Chat.ID, fmt.Sprintf(res.TxtAddSongSuccess, songName))
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		logrus.WithError(err).Error("Unable to send add song response")
	}
}
