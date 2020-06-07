package command

import (
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type voteSkip struct {
	stateSvc *state.Service
	bot      *bot.BotAPI
}

func (s *voteSkip) accepts(update bot.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == res.CmdVoteSkip
}

func (s *voteSkip) Handle(update bot.Update) () {
	_, _ = s.bot.DeleteMessage(bot.DeleteMessageConfig{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.MessageID,
	})
	num, err := s.bot.GetChatMembersCount(update.Message.Chat.ChatConfig())
	if err != nil {
		logrus.WithError(err).Error("Unable to get num of chat members")
		return
	}
	resp, err := s.stateSvc.SkipSong(update.Message.From, update.Message.Chat, num)
	var msg bot.MessageConfig
	if err != nil {
		msg = bot.NewMessage(update.Message.Chat.ID, err.Error())
	} else {
		msg = bot.NewMessage(update.Message.Chat.ID, resp)
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		logrus.WithError(err).Error("Unable to send vote skip response")
	}
}
