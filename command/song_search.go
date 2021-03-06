package command

import (
	"fmt"
	app "github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/res"
	"github.com/arttor/spoty-paty-bot/state"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"strings"
)

const (
	searchCallbackPrefix = "SPB_SEARCH:"
)

type songSearch struct {
	stateSvc *state.Service
	bot      *bot.BotAPI
}

func (s *songSearch) Handle(update bot.Update) () {
	if update.CallbackQuery != nil {
		s.handleCallback(update)
	} else {
		s.handleCommand(update)
	}
}
func (s *songSearch) accepts(update bot.Update) bool {
	return (update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, searchCallbackPrefix)) || (update.Message != nil && update.Message.IsCommand() && update.Message.Command() == res.CmdSearch)
}

func (s *songSearch) handleCallback(update bot.Update) {
	songID := strings.TrimPrefix(update.CallbackQuery.Data, searchCallbackPrefix)
	chat, ok := s.stateSvc.Get(update.CallbackQuery.Message.Chat.ID)
	if !ok || chat.DjID == 0 {
		_, _ = s.bot.AnswerCallbackQuery(bot.NewCallback(update.CallbackQuery.ID, res.TxtSearchSongNoDj))
	}
	err := s.stateSvc.AddSong(update.CallbackQuery.From, update.CallbackQuery.Message.Chat, spotify.ID(songID))
	if err == nil {
		_, err = s.bot.AnswerCallbackQuery(bot.NewCallback(update.CallbackQuery.ID, res.TxtCallbackAddSongSuccess))
		if err != nil {
			logrus.WithError(err).Error("Unable to send search callback response")
		}
		return
	}
	if strings.Contains(err.Error(), res.TxtNoDjError) || strings.HasSuffix(err.Error(), "DJ can use /settings command to increase max songs number.") {
		_, err = s.bot.Send(bot.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error()))
		if err != nil {
			logrus.WithError(err).Error("Unable handle search callback")
		}
		return
	}
	_, err = s.bot.AnswerCallbackQuery(bot.NewCallback(update.CallbackQuery.ID, err.Error()))
	if err != nil {
		logrus.WithError(err).Error("Unable to send search callback response")
	}
}

func (s *songSearch) handleCommand(update bot.Update) {
	searchQuery := update.Message.CommandArguments()
	if searchQuery == "" {
		_, err := s.bot.Send(bot.NewMessage(update.Message.Chat.ID, res.TxtSearchSongEmptyQuery))
		if err != nil {
			logrus.WithError(err).Error("Unable handle search callback")
		}
		return
	}
	songs, err := s.stateSvc.SearchSongs(update.Message.Chat, searchQuery)
	if err != nil {
		_, err = s.bot.Send(bot.NewMessage(update.Message.Chat.ID, err.Error()))
		if err != nil {
			logrus.WithError(err).Error("Unable handle search callback")
		}
		return
	}
	response := bot.NewMessage(update.Message.Chat.ID, fmt.Sprintf(res.TxtSearchResultPattern, len(songs), searchQuery))
	btns := make([][]bot.InlineKeyboardButton, len(songs))
	for i, song := range songs {
		btns[i] = bot.NewInlineKeyboardRow(bot.NewInlineKeyboardButtonData(songPresentation(song), searchCallbackPrefix+string(song.ID)),
		)
	}
	response.ReplyMarkup = bot.NewInlineKeyboardMarkup(btns...)
	_, err = s.bot.Send(response)
	if err != nil {
		logrus.WithError(err).Error("Unable to send search response")
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
