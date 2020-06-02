package state

import "github.com/zmb3/spotify"

type Chat struct {
	Id            int64
	MaxSongs      int
	Songs         int
	SongsUserID   int
	SpotifyToken  string
	SpotifyClient *spotify.Client
}
