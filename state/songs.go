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

func (s *Service) AddSong(fromUser *bot.User, fromChat *bot.Chat, songID spotify.ID) error {
	s.m.Lock()
	defer s.m.Unlock()
	chatID := fromChat.ID
	user := fromUser.String()
	chat, ok := s.mem[chatID]
	if !ok || chat.DjID == 0 {
		return errors.New(res.TxtNoDjError)
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
		return errors.New(res.TxtNoActiveDeviceError)
	}
	s.mem[chatID] = chat
	return nil
}

func (s *Service) SearchSongs(fromChat *bot.Chat, searchQuery string) ([]spotify.FullTrack, error) {
	s.m.RLock()
	chat, ok := s.mem[fromChat.ID]
	s.m.RUnlock()
	if !ok || chat.DjID == 0 {
		return nil, errors.New(res.TxtNoDjError)
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

func (s *Service) GetQueue(fromChat *bot.Chat) ([]spotify.FullTrack, error) {
	s.m.RLock()
	chat, ok := s.mem[fromChat.ID]
	s.m.RUnlock()
	if !ok || chat.DjID == 0 {
		return nil, errors.New(res.TxtNoDjError)
	}
	// TODO: wait spotify get queue api or maintain queue
	return nil, errors.New("not implemented")
}

func (s *Service) SkipSong(fromUser *bot.User, fromChat *bot.Chat, numOfChatMembers int) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()
	chat, ok := s.mem[fromChat.ID]
	if !ok || chat.DjID == 0 {
		return "", errors.New(res.TxtNoDjError)
	}
	client := chat.DjClient
	state, err := client.PlayerState()
	if err != nil {
		logrus.WithError(err).Error("Spotify get state error")
		return "", errors.New(res.TxtNoActiveDeviceError)
	}
	currSong := state.CurrentlyPlaying.Item
	if currSong == nil {
		return "", errors.New(res.TxtNoSongError)
	}
	if chat.VoteSkip.SongID != currSong.ID {
		chat.VoteSkip.SongID = currSong.ID
		chat.VoteSkip.SongName = currSong.Name
		chat.VoteSkip.Votes = nil
	}
	if chat.VoteSkip.Votes == nil {
		chat.VoteSkip.Votes = make(map[int]struct{})
	}
	_, isVoted := chat.VoteSkip.Votes[fromUser.ID]
	votesNeeded := numOfChatMembers / 2
	if isVoted {
		s.mem[fromChat.ID] = chat
		return "", errors.New(fmt.Sprintf(res.TxtVoteSkipAlreadyVotedPattern, fromUser.String(), len(chat.VoteSkip.Votes), votesNeeded, currSong.Name))
	}
	chat.VoteSkip.Votes[fromUser.ID] = struct{}{}
	s.mem[fromChat.ID] = chat
	votes := len(chat.VoteSkip.Votes)
	if votes >= votesNeeded {
		err = chat.DjClient.Next()
		if err != nil {
			logrus.WithError(err).Error("Spotify get state error")
			return "", errors.New(res.TxtNoActiveDeviceError)
		}
		return fmt.Sprintf(res.TxtVoteSkipSuccessPattern, len(chat.VoteSkip.Votes), votesNeeded, currSong.Name), nil
	} else {
		return fmt.Sprintf(res.TxtVoteSkipVotedPattern, len(chat.VoteSkip.Votes), votesNeeded, currSong.Name), nil
	}
}
