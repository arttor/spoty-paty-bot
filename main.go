package main

import (
	"github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/service"
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

func main() {
	logrus.Info("Starting bot app...")
	conf, err := app.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	//db := setupDB()
	//defer func() {
	//	logrus.Error(db.Close())
	//}()
	stateSvc:=state.New(conf,nil)
	spotifySvc:=spotify.New(conf,stateSvc)
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to register bot with token")
	}
	router:=service.New(stateSvc,spotifySvc,bot)
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
		if update.Message != nil {
			log.Printf("msg: %+v\n", *update.Message)
		}
		router.Handle(update)
	}
}