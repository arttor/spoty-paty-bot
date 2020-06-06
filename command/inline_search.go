package command

import (
	"fmt"
	"github.com/arttor/spoty-paty-bot/search"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"strconv"
	"strings"
)

var (
	limit = 10
)

type inlineSearch struct {
	searchSvc search.Service
	bot      *bot.BotAPI
}

func (s *inlineSearch) accepts(update bot.Update) bool {
	return update.InlineQuery != nil
}

func (s *inlineSearch) Handle(update bot.Update) () {
	query := update.InlineQuery
	client := s.searchSvc.GetClient()
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
		id := fmt.Sprintf("sppbid:%s:69", track.ID)
		artist := ""
		for _, a := range track.Artists {
			artist = artist + a.Name + ", "
		}
		artist = strings.TrimSuffix(artist, ", ")
		r := bot.NewInlineQueryResultArticle(id, getSongPresentation(track), "")
		r.Description = artist
		if len(track.Album.Images) > 0 {
			r.ThumbURL = track.Album.Images[0].URL
			r.ThumbHeight = track.Album.Images[0].Height
			r.ThumbWidth = track.Album.Images[0].Width
		}
		r.InputMessageContent = bot.InputTextMessageContent{
			Text: "/search " + track.Name,
		}
		results[i] = r
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

func getSongPresentation(song spotify.FullTrack) string {
	sec := song.Duration / 1000
	return fmt.Sprintf("[%v:%02d] %s", sec/60, sec%60, song.Name)
}