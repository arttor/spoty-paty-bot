package service

import (
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type login struct {
	next       Handler
	stateSvc   *state.Service
	spotifySvc *spotify.Service
	bot        *bot.BotAPI
}

func (s *login) Handle(update tgbotapi.Update) () {
	if s.accepts(update) {
		if update.Message.IsCommand() {
			s.handleCommand(update)
		} else {
			s.handleCallback(update)
		}
		return
	}
	if s.next != nil {
		s.next.Handle(update)
		return
	}
	logrus.Info("No handler for given update")
}
func (s *login) accepts(update tgbotapi.Update) bool {
	return (update.Message.IsCommand() && update.Message.Command() == "login") || (update.CallbackQuery != nil )
}

func (s *login) handleCommand(update tgbotapi.Update) {
	chat, err := s.stateSvc.Get(update.Message.Chat.ID)
	if err != nil || chat.SpotifyToken == "" {
		s.login(update)
	} else {
		s.alreadyLoggedIn(update)
	}
}

func (s *login) handleCallback(update tgbotapi.Update) {
	logrus.Infof("CALLBACK LOGIN FROM %s", update.CallbackQuery.From.UserName)
	resp, err := s.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, update.CallbackQuery.Data))
	if err!=nil{
		logrus.WithError(err).Error("Login Callback response erro")
	}else {
		logrus.Infof("CALBACK login response of response %v",resp.Description)
	}
}

func (s *login) login(update tgbotapi.Update) {
	url := s.spotifySvc.GetAuthURL(update.Message.Chat.ID)
	response := tgbotapi.NewMessage(update.Message.Chat.ID, "Please grant access to play songs")
	btn := tgbotapi.NewInlineKeyboardButtonURL("Login into Spotify account", url)
	callbackData := "login"
	btn.CallbackData = &callbackData
	response.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		btn,
	))
	_, err := s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send log in request")
	}
}

func (s *login) alreadyLoggedIn(update tgbotapi.Update) {
	response := tgbotapi.NewMessage(update.Message.Chat.ID, "Already logged in")
	_, err := s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send already logged in response")
	}
}
