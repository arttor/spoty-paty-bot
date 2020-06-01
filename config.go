package main

import (
	"errors"
	"flag"
	"os"
)

const (
	serverPort      = "8080"
	tokenFlagName   = "token"
	tokenEnvName    = "TG_BOT_TOKEN"
	webHookFlagName = "wh"
	webHookEnvName  = "TG_BOT_WEB_HOOK"
)

type config struct {
	webHookBaseURL string
	token          string
}

func readConfig() (config, error) {
	res := config{}
	host := flag.String(webHookFlagName, "", "web hook base url")
	token := flag.String(tokenFlagName, "", "telegram bot token")
	flag.Parse()
	if host != nil && *host != "" {
		res.webHookBaseURL = *host
	} else {
		res.webHookBaseURL = os.Getenv(tokenEnvName)
	}
	if res.webHookBaseURL == "" {
		return res, errors.New("web hook base url not specified")
	}
	if token != nil && *token != "" {
		res.token = *token
	} else {
		res.token = os.Getenv(webHookEnvName)
	}
	if res.token == "" {
		return res, errors.New("telegram token not specified")
	}
	return res, nil
}
