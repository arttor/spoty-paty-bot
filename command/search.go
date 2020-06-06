package command

import (
	"fmt"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

var (
	limit = 10
)

type search struct {
	bot      *bot.BotAPI
}

func (s *search) accepts(update bot.Update) bool {
	return update.InlineQuery != nil
}

func (s *search) Handle(update bot.Update) () {
	query := update.InlineQuery
	if len([]rune(query.Query)) < 4 {
		return
	}
	results := make([]interface{}, 3)
	for i:=0;i<3;i++ {
		r := bot.NewInlineQueryResultArticle(fmt.Sprintf("id%v",i), fmt.Sprintf("songname%v",i), "")
		r.Description = "artist"
		//if len(track.Album.Images) > 0 {
		//	r.ThumbURL = track.Album.Images[0].URL
		//	r.ThumbHeight = track.Album.Images[0].Height
		//	r.ThumbWidth = track.Album.Images[0].Width
		//}
		r.InputMessageContent = bot.InputTextMessageContent{
			Text: "/search писька",
		}
		results[i] = r
	}
	_, err := s.bot.AnswerInlineQuery(bot.InlineConfig{
		InlineQueryID: query.ID,
		Results:       results,
		CacheTime:     600,
		IsPersonal:    false,
		NextOffset:    "",
	})
	if err != nil {
		logrus.WithError(err).Error("Unable to send logout response")
	}
}
