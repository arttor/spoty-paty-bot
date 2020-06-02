package service

import (
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func New(stateSvc *state.Service, spotifySvc *spotify.Service, bot *bot.BotAPI) Handler {
	return &login{
		next: &logout{
			next:     nil,
			stateSvc: stateSvc,
			bot:      bot,
		},
		stateSvc:   stateSvc,
		spotifySvc: spotifySvc,
		bot:        bot,
	}
}
