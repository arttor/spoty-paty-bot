package main

import (
	"context"
	"github.com/arttor/spoty-paty-bot/command"
	"github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/db"
	"github.com/arttor/spoty-paty-bot/search"
	"github.com/arttor/spoty-paty-bot/spotify"
	"github.com/arttor/spoty-paty-bot/state"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	database, err := db.New()
	if err != nil {
		logrus.WithError(err).Fatal("Inline search bot Unable to create db")
	}
	searchSvc := search.NewService(database)
	spotifySvc := spotify.New(conf, stateSvc, bot, searchSvc)
	searchSvc.RestoreClient(spotifySvc.GetClient)
	logrus.Infof("SearchLoginURL: %v", spotifySvc.GetSearchAuthURL())
	router := command.New(stateSvc, spotifySvc, bot, searchSvc)
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
	ctx := listenSignal()
	time.Sleep(time.Millisecond * 500)
	updates.Clear()
	logrus.Info("Listening for updates...")
	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				log.Printf("msg: %+v\n", *update.Message)
			}
			router.Handle(update)
		case <-ctx.Done():
			logrus.Info("Gracefully stopping")
			searchSvc.Close()
			err = database.Close()
			if err != nil {
				logrus.WithError(err).Error("Inline search bot Unable to close db")
			}
			logrus.Info("Gracefully stopped")
		}
	}
}

func listenSignal() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()

	}()
	return ctx
}
