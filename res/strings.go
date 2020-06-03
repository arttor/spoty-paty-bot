package res

// BOT TEXT MESSAGES
const (
	TxtLoginInfo                        = "Please grant access to play songs"
	TxtLoginBtn                         = "Login into Spotify account"
	TxtLoginAlreadyPattern              = "This party already have DJ: @%s \nPS /kick and /logout commands may help"
	TxtLogoutErrNotLogin                = "You need to login first. Use /login command"
	TxtLogoutErrAnotherUserPattern      = "Cannot logout. You are not DJ like @%s\n PS see /kick command"
	TxtLogoutSuccess                    = "Logout Success. Party has no DJ..."
	TxtFinishLoginSuccessPattern        = "Welcome! @%s is a DJ"
	TxtAddSongNoDj                      = "Party has no DJ. Someone should /login"
	TxtAddSongToMuchSongsInQueuePattern = "Easy, @%s! It will be %v song in a row from you. DJ can use /settings command to increase max songs number."
	TxtAddSongSpotifyError              = "Spotify queue song error"
	TxtAddSongSuccess                   = "Song added to queue. Nice choice!"
	TxtCallbackAddSongSuccess           = "Song added!"
	TxtSearchSongSpotifyError           = "Spotify search song error"
	TxtSearchSongNoSongsFoundError      = "No such song"
	TxtSearchResultPattern              = "Top %v songs for '%s' found:"
	TxtSearchSongEmptyQuery             = "Song name not specified"
)

// WEB TEXT MESSAGES
const (
	TxtWebLoginAlready    = "Too late! This party already have DJ"
	TxtWebLoginNoSuchChat = "no such chat"
)

// BOT COMMANDS
const (
	CmdLogin       = "login"
	CmdLoginFinish = "becameDj"
	CmdLogout      = "logout"
	CmdSearch      = "search"
)
