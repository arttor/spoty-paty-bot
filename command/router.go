package command

import (
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	Handle(update bot.Update)
	accepts(update bot.Update) bool
}

type router struct {
	handlers []Handler
}

func New(stateSvc *state.Service, spotifySvc *spotify.Service, bot *bot.BotAPI) Handler {
	return &router{handlers: []Handler{
		&login{
			stateSvc:   stateSvc,
			spotifySvc: spotifySvc,
			bot:        bot},
		&logout{
			stateSvc: stateSvc,
			bot:      bot},
		&loginFinish{
			stateSvc: stateSvc,
			bot:      bot},
		&songLink{
			stateSvc: stateSvc,
			bot:      bot},
		&songSearch{
			stateSvc: stateSvc,
			bot:      bot},
		&djLeftChat{
			stateSvc: stateSvc,
			bot:      bot},
	}}
}

func (r *router) Handle(update bot.Update) {
	for _, handler := range r.handlers {
		if handler.accepts(update) {
			handler.Handle(update)
			return
		}
	}
	logrus.Info("No handler for given update")
}

func (r *router) accepts(update bot.Update) bool {
	return true
}
