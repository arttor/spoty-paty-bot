package spotify

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (s *Service) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["state"]
	if !ok || len(keys[0]) < 1 {
		logrus.Error("Url Param 'state' is missing")
		http.Error(w, "Url Param 'state' is missing", http.StatusBadRequest)
		return
	}
	fmt.Printf("------------ %+v", keys)
	state := keys[0]
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
	client := s.auth.NewClient(token)
	err = s.state.SaveClient(n, token, &client)
	if err == nil {
		logrus.WithError(err).Error("Unable to save token")
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "Authenticated!")
}
