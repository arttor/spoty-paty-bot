package state

import "github.com/zmb3/spotify"

type Chat struct {
	Id              int64
	MaxSongs        int
	DjClient        *spotify.Client
	DjID            int
	DjName          string
	LoginCandidates map[string]*spotify.Client
	Queue           map[string]int
	VoteSkip        VoteSkip
}
type VoteSkip struct {
	SongID   spotify.ID
	SongName string
	Votes    map[int]struct{}
}
