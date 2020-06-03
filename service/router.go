package service

import (
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func New(stateSvc *state.Service, spotifySvc *spotify.Service, bot *bot.BotAPI) Handler {
	return &login{
		stateSvc:   stateSvc,
		spotifySvc: spotifySvc,
		bot:        bot,
		next: &logout{
			stateSvc: stateSvc,
			bot:      bot,
			next: &loginFinish{
				stateSvc: stateSvc,
				bot:      bot,
				next:     &songLink{
					stateSvc: stateSvc,
					bot:      bot,
					next:     &songSearch{
						stateSvc: stateSvc,
						bot:      bot,
						next:     nil,
					},
				},
			},
		},
	}
}
