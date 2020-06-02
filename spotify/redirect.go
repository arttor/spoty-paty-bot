package spotify

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
)

const (
	FinishLoginCommandPattern = "/" + FinishLoginCommand + " %s"
	FinishLoginCommand        = "i_am_DJ_now"
)

func (s *Service) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["state"]
	if !ok || len(keys[0]) < 1 {
		logrus.Error("Url Param 'state' is missing")
		http.Error(w, "Url Param 'state' is missing", http.StatusBadRequest)
		return
	}
	state := keys[0]
	token, err := s.auth.Token(state, r)
	if err != nil {
		logrus.WithError(err).Error("Unable to get access token from code")
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	n, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		logrus.WithError(err).Error("Unable to convert state")
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	client := s.auth.NewClient(token)
	code, err := s.state.SaveClient(n, &client)
	if err != nil {
		logrus.WithError(err).Error("Unable to save token")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("res/login.html")
	if err != nil {
		logrus.WithError(err).Error("Error parsing login template")
		http.Error(w, "Login error", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, ViewData{LoginCommand: fmt.Sprintf(FinishLoginCommandPattern, code)})
}

type ViewData struct {
	LoginCommand string
}
