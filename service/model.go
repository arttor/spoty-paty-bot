package service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Handler interface {
	Handle(update tgbotapi.Update)
	Accepts(update tgbotapi.Update) bool
}
