package spotify

import (
	app "github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/search"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"os"
	"strconv"
)

const (
	Callback = "/spotify/callback/"
)

type Service struct {
	state           *state.Service
	search          search.Service
	config          app.Config
	auth            *spotify.Authenticator
	searchAuth      *spotify.Authenticator
	bot             *bot.BotAPI
	searchAuthState string
}

func New(config app.Config, state *state.Service, bot *bot.BotAPI, search search.Service) *Service {
	auth := spotify.NewAuthenticator(config.BaseURL+Callback, spotify.ScopeStreaming, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	searchAuth := spotify.NewAuthenticator(config.BaseURL + Callback)
	auth.SetAuthInfo(config.SpotifyClientID, config.SpotifyClientSecret)
	searchAuth.SetAuthInfo(config.SpotifyClientID, config.SpotifyClientSecret)
	return &Service{config: config, state: state, auth: &auth, searchAuth: &searchAuth, bot: bot, search: search, searchAuthState: os.Getenv("SEARCH_STATE")}
}

func (s *Service) GetAuthURL(chatID int64) string {
	s.state.AddIfNotExists(state.Chat{
		Id:       chatID,
		MaxSongs: app.DefaultMaxSongs,
	})
	return s.auth.AuthURLWithDialog(strconv.FormatInt(chatID, 10))
}

func (s *Service) GetSearchAuthURL() string {
	return s.searchAuth.AuthURLWithDialog(s.searchAuthState)
}


func (s *Service) GetClient(token *oauth2.Token) spotify.Client {
	return s.searchAuth.NewClient(token)
}