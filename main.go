package main

import (
	"context"
	"github.com/arttor/spoty-paty-bot/command"
	"github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/inlinesearch"
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	logrus.Info("Starting bot app...")
	conf, err := app.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to register bot with token")
	}
	stateSvc := state.New(conf, nil)
	spotifySvc := spotify.New(conf, stateSvc, bot)
	router := command.New(stateSvc, spotifySvc, bot)
	logrus.Info("All services started")
	app.SetupLog(bot)
	logrus.Infof("Authorized on account %s", bot.Self.UserName)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(conf.BaseURL + "/" + bot.Token))
	if err != nil {
		logrus.WithError(err).Fatal("Unable to register web-hook")
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		logrus.WithError(err).Fatal("Unable to ge web-hook info")
	}
	if info.LastErrorDate != 0 {
		logrus.Warnf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	http.HandleFunc(spotify.Callback, spotifySvc.RedirectHandler)
	go func() {
		logrus.Error(http.ListenAndServe("0.0.0.0:"+conf.Port, nil))
	}()
	go func() {
		logrus.Error(inlinesearch.Start(context.Background()))
	}()
	time.Sleep(time.Millisecond * 500)
	updates.Clear()
	logrus.Info("Listening for updates...")
	for update := range updates {
		log.Printf("%+v\n", update)
		if update.Message != nil {
			log.Printf("msg: %+v\n", *update.Message)
		}
		router.Handle(update)
	}
}
