package app

import (
	"errors"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	defaultPort = "8080"
	logLvl      = "info"
	DefaultMaxSongs      = 3
)

type Config struct {
	BaseURL             string
	Token               string
	Port                string
	SpotifyClientID     string
	SpotifyClientSecret string
}

func ReadConfig() (Config, error) {
	res := Config{}
	res.Port = os.Getenv("PORT")
	if res.Port == "" {
		res.Port = defaultPort
	}
	res.BaseURL = os.Getenv("TG_BOT_URL")
	if res.BaseURL == "" {
		return res, errors.New("web hook base url not specified")
	}
	res.Token = os.Getenv("TG_BOT_TOKEN")
	if res.Token == "" {
		return res, errors.New("telegram Token not specified")
	}
	res.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	if res.SpotifyClientID == "" {
		return res, errors.New("SpotifyClientID not specified")
	}
	res.SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	if res.SpotifyClientSecret == "" {
		return res, errors.New("SpotifyClientSecret not specified")
	}
	return res, nil
}

func SetupLog(api *bot.BotAPI) {
	switch logLvl {
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
		api.Debug = false
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		api.Debug = false
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
		api.Debug = false
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		api.Debug = false
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		api.Debug = true
	}
}
