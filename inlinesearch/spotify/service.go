package spotify

import (
	"fmt"
	"github.com/arttor/spoty-paty-bot/inlinesearch/client"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

const (
	Callback                  = "/spotify/search/callback/"
)

type Service struct {
	auth      *spotify.Authenticator
	authState string
	search    client.Service
}

func New(search client.Service) *Service {
	auth := spotify.NewAuthenticator(os.Getenv("TG_BOT_URL") + Callback)
	auth.SetAuthInfo(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	return &Service{auth: &auth, authState: os.Getenv("SEARCH_STATE"), search: search}
}

func (s *Service) GetAuthURL() string {
	return s.auth.AuthURL(s.authState)
}

func (s *Service) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["state"]
	if !ok || len(keys[0]) < 1 {
		logrus.Error("Url Param 'state' is missing")
		http.Error(w, "Url Param 'state' is missing", http.StatusBadRequest)
		return
	}
	state := keys[0]
	if state != s.authState {
		http.Error(w, "invalid state", http.StatusNotFound)
		return
	}
	token, err := s.auth.Token(state, r)
	if err != nil {
		logrus.WithError(err).Error("Unable to get access token from code")
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	client := s.auth.NewClient(token)
	s.search.SetClient(client)
	w.WriteHeader(200)
	_, _ = fmt.Fprint(w, "Authenticated")
}


func (s *Service) GetClient(token *oauth2.Token) spotify.Client {
	return s.auth.NewClient(token)
}