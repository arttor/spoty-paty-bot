package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func main() {
	conf, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.token)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to register bot with token")
	}
	logrus.SetLevel(logrus.DebugLevel)
	bot.Debug = true
	logrus.Printf("Authorized on account %s", bot.Self.UserName)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(conf.webHookBaseURL + "/" + bot.Token))
	if err != nil {
		logrus.WithError(err).Fatal("Unable to register webhook")
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		logrus.WithError(err).Fatal("Unable to ge webhook info")
	}
	if info.LastErrorDate != 0 {
		logrus.Warnf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	go func() {
		logrus.Error(http.ListenAndServe("0.0.0.0:"+conf.port, nil))
	}()
	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
