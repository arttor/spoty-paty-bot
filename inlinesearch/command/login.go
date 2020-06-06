package command

import (
	"github.com/arttor/spoty-paty-bot/inlinesearch"
	"github.com/arttor/spoty-paty-bot/inlinesearch/spotify"
	"github.com/arttor/spoty-paty-bot/res"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type login struct {
	stateSvc   inlinesearch.Service
	spotifySvc *spotify.Service
	bot        *bot.BotAPI
}

func (s *login) Handle(update tgbotapi.Update) () {
	client:=s.stateSvc.GetClient()
	if client == nil {
		s.login(update)
	} else {
		s.alreadyLoggedIn(update)
	}
}

func (s *login) accepts(update tgbotapi.Update) bool {
	return update.Message!=nil && update.Message.IsCommand() && update.Message.Command() == res.CmdLogin
}

func (s *login) login(update tgbotapi.Update) {
	url := s.spotifySvc.GetAuthURL()
	response := tgbotapi.NewMessage(update.Message.Chat.ID, res.TxtLoginInfo)
	response.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL(res.TxtLoginBtn, url),
	))
	_, err := s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send log in request")
	}
}

func (s *login) alreadyLoggedIn(update tgbotapi.Update) {
	_, err := s.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Already logged. use /logout"))
	if err != nil {
		logrus.WithError(err).Error("Unable to send already logged in response")
	}
}
