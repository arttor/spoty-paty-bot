package state

import (
	"errors"
	"fmt"
	"github.com/arttor/spoty-paty-bot/res"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
)

func (s *Service) QueueSong(update bot.Update, songID spotify.ID) error {
	s.m.Lock()
	defer s.m.Unlock()
	chatID := update.Message.Chat.ID
	user := update.Message.From.String()
	chat, ok := s.mem[chatID]
	if !ok || chat.DjID == 0 {
		return errors.New(res.TxtAddSongNoDj)
	}
	userSongsInARow := chat.Queue[user]
	if userSongsInARow >= chat.MaxSongs {
		return errors.New(fmt.Sprintf(res.TxtAddSongToMuchSongsInQueuePattern, user, chat.MaxSongs))
	}
	chat.Queue = make(map[string]int)
	userSongsInARow++
	chat.Queue[user] = userSongsInARow
	err:=chat.DjClient.QueueSong(songID)
	if err != nil {
		logrus.WithError(err).Error("Spotify queue song error")
		return errors.New("Spotify queue song error")
	}
	return nil
}
