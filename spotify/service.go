package spotify

import (
	"fmt"
	app "github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/state"
	"github.com/zmb3/spotify"
	"strconv"
)

const (
	Callback = "/spotify/callback/"
)

type Service struct {
	state  *state.Service
	config app.Config
	auth   *spotify.Authenticator
}

func New(config app.Config, state *state.Service) *Service {
	auth := spotify.NewAuthenticator(config.BaseURL+Callback, spotify.ScopeStreaming, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	auth.SetAuthInfo(config.SpotifyClientID, config.SpotifyClientSecret)
	return &Service{config: config, state: state, auth: &auth}
}

func (s *Service) GetAuthURL(chatID int64) string {
	_ = s.state.AddIfExists(state.Chat{
		Id:       chatID,
		MaxSongs: app.DefaultMaxSongs,
	})
	fmt.Printf("==========%v --%v",strconv.FormatInt(chatID, 10),chatID)
	return s.auth.AuthURLWithDialog(strconv.FormatInt(chatID, 10))
}
