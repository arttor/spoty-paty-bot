package main

import (
	"github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/service"
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
	router := service.New(stateSvc, spotifySvc, bot)
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
	time.Sleep(time.Millisecond * 500)
	updates.Clear()
	logrus.Info("Listening for updates...")
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("%+v\n", update)
		log.Printf("------------- username %v fm %v lm%v\n", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName)
		if update.Message != nil {
			log.Printf("msg: %+v\n", *update.Message)
		}
		router.Handle(update)
	}
}
