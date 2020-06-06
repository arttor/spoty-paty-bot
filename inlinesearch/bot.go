package inlinesearch

import (
	"context"
	"github.com/arttor/spoty-paty-bot/db"
	"github.com/arttor/spoty-paty-bot/inlinesearch/client"
	"github.com/arttor/spoty-paty-bot/inlinesearch/command"
	"github.com/arttor/spoty-paty-bot/inlinesearch/spotify"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const (
	envBotTokenName = "TG_SEARCH_BOT_TOKEN"
)

func Start(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(envBotTokenName))
	if err != nil {
		logrus.WithError(err).Error("Unable to register bot with token")
		return err
	}
	logrus.Infof("Inline search bot Authorized on account %s", bot.Self.UserName)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("TG_BOT_URL") + "/" + bot.Token))
	if err != nil {
		logrus.WithError(err).Error("Inline search botUnable to register web-hook")
		return err
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		logrus.WithError(err).Error("Inline search bot Unable to ge web-hook info")
		return err
	}
	if info.LastErrorDate != 0 {
		logrus.Warnf("Telegram inline search callback failed: %s", info.LastErrorMessage)
	}
	database, err := db.New()
	if err != nil {
		logrus.WithError(err).Error("Inline search bot Unable to create db")
		return err
	}
	svc := client.NewService(database)
	spotifySvc := spotify.New(svc)
	svc.RestoreClient(spotifySvc.GetClient)
	http.HandleFunc(spotify.Callback, spotifySvc.RedirectHandler)
	updates := bot.ListenForWebhook("/" + bot.Token)
	handler := command.New(svc, spotifySvc, bot)
	time.Sleep(time.Millisecond * 500)
	updates.Clear()
	logrus.Info("Inline search Listening for updates...")
	for {
		select {
		case update := <-updates:
			handler.Handle(update)
		case <-ctx.Done():
			svc.Close()
			err = database.Close()
			if err != nil {
				logrus.WithError(err).Error("Inline search bot Unable to close db")
			}
			return nil
		}
	}
}
