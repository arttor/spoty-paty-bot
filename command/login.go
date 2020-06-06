package command

import (
	"fmt"
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type login struct {
	stateSvc   *state.Service
	spotifySvc *spotify.Service
	bot        *bot.BotAPI
}

func (s *login) Handle(update tgbotapi.Update) () {
	chat, _ := s.stateSvc.Get(update.Message.Chat.ID)
	if chat.DjID == 0 {
		s.login(update)
	} else {
		s.alreadyLoggedIn(update, chat)
	}
}
func (s *login) accepts(update tgbotapi.Update) bool {
	return update.Message!=nil && update.Message.IsCommand() && update.Message.Command() == res.CmdLogin
}

func (s *login) login(update tgbotapi.Update) {
	url := s.spotifySvc.GetAuthURL(update.Message.Chat.ID)
	response := tgbotapi.NewMessage(update.Message.Chat.ID, res.TxtLoginInfo)
	response.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL(res.TxtLoginBtn, url),
	))
	_, err := s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send log in request")
	}
}

func (s *login) alreadyLoggedIn(update tgbotapi.Update, chat state.Chat) {
	response := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(res.TxtLoginAlreadyPattern, chat.DjName))
	_, err := s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send already logged in response")
	}
}
