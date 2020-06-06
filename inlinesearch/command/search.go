package command

import (
	"fmt"
	app "github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/inlinesearch/client"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"strconv"
	"strings"
)

var (
	limit = 10
)

type search struct {
	stateSvc client.Service
	bot      *bot.BotAPI
}

func (s *search) accepts(update bot.Update) bool {
	return update.InlineQuery != nil
}

func (s *search) Handle(update bot.Update) () {
	query := update.InlineQuery
	client := s.stateSvc.GetClient()
	if client == nil {
		return
	}
	if len([]rune(query.Query)) < 4 {
		return
	}
	offset, _ := strconv.Atoi(query.Offset)
	res, err := client.SearchOpt(query.Query, spotify.SearchTypeTrack, &spotify.Options{Offset: &offset, Limit: &limit})
	if err != nil {
		logrus.WithError(err).Error("Spotify inline search error")
		return
	}
	if res == nil || res.Tracks == nil {
		logrus.WithError(err).Error("Nil result")
		return
	}
	nextOffsetInt := offset + limit
	nextOffsetStr := strconv.Itoa(nextOffsetInt)
	if nextOffsetInt >= res.Tracks.Total {
		nextOffsetStr = ""
	}
	results := make([]interface{}, len(res.Tracks.Tracks))
	for i, track := range res.Tracks.Tracks {
		results[i] = bot.NewInlineQueryResultAudio(string(track.ID), track.PreviewURL, songPresentation(track))
	}
	_, err = s.bot.AnswerInlineQuery(bot.InlineConfig{
		InlineQueryID: query.ID,
		Results:       results,
		CacheTime:     600,
		IsPersonal:    false,
		NextOffset:    nextOffsetStr,
	})
	if err != nil {
		logrus.WithError(err).Error("Unable to send logout response")
	}
}

func songPresentation(song spotify.FullTrack) string {
	artist := ""
	for _, a := range song.Artists {
		artist = artist + a.Name + ", "
	}
	artist = strings.TrimSuffix(artist, ", ")
	if len([]rune(artist)) > app.SongSearchMaxArtistLength {
		artist = string([]rune(artist)[:app.SongSearchMaxArtistLength-3]) + "..."
	}
	songName := song.Name
	if len([]rune(songName)) > app.SongSearchMaxSongLength {
		songName = string([]rune(songName)[:app.SongSearchMaxSongLength-3]) + "..."
	}
	sec := song.Duration / 1000
	return fmt.Sprintf("%s - %s   %v:%02d", songName, artist, sec/60, sec%60)
}
