package state

import (
	"errors"
	"fmt"
	app "github.com/arttor/spoty-paty-bot/config"
	"github.com/arttor/spoty-paty-bot/res"
	bot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
)

func (s *Service) QueueSong(fromUser *bot.User, fromChat *bot.Chat, songID spotify.ID) error {
	s.m.Lock()
	defer s.m.Unlock()
	chatID := fromChat.ID
	user := fromUser.String()
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
	err := chat.DjClient.QueueSong(songID)
	if err != nil {
		logrus.WithError(err).Error("Spotify queue song error")
		return err
	}
	s.mem[chatID]=chat
	return nil
}

func (s *Service) SearchSongs(fromChat *bot.Chat, searchQuery string) ([]spotify.FullTrack, error) {
	s.m.RLock()
	chat, ok := s.mem[fromChat.ID]
	s.m.RUnlock()
	if !ok || chat.DjID == 0 {
		return nil, errors.New(res.TxtAddSongNoDj)
	}
	client := chat.DjClient
	result, err := client.Search(searchQuery, spotify.SearchTypeTrack)
	if err != nil {
		logrus.WithError(err).Error("Spotify search song error")
		return nil, errors.New(res.TxtSearchSongSpotifyError)
	}
	if result.Tracks == nil || result.Tracks.Tracks == nil || len(result.Tracks.Tracks) == 0 {
		return nil, errors.New(res.TxtSearchSongNoSongsFoundError)
	}
	songs := make([]spotify.FullTrack, 0)
	for i := 0; i < len(result.Tracks.Tracks) && i < app.SongSearchMaxResult; i++ {
		songs = append(songs, result.Tracks.Tracks[i])
	}
	return songs, nil
}
