package spotify

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (s *Service) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	// use the same state string here that you used to generate the URL
	token, err := s.auth.Token(state, r)
	if err != nil {
		logrus.WithError(err).Error("Unable to get access token from code")
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	n, err := strconv.ParseInt(state, 10, 64)
	if err == nil {
		logrus.WithError(err).Error("Unable to convert state")
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	err = s.state.SaveToken(n, token)
	if err == nil {
		logrus.WithError(err).Error("Unable to save token")
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	w.WriteHeader(200)
}
