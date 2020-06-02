package res

// BOT TEXT MESSAGES
const (
	TxtLoginInfo    = "Please grant access to play songs"
	TxtLoginBtn     = "Login into Spotify account"
	TxtLoginAlreadyPattern = "This party already have DJ: @%s \nPS /kick and /logout commands may help"
	TxtLogoutErrNotLogin = "You need to login first. Use /login command"
	TxtLogoutErrAnotherUserPattern = "Cannot logout. You are not DJ like @%s\n PS see /kick command"
	TxtLogoutSuccess = "Logout Success. Party has no DJ..."
	TxtFinishLoginSuccessPattern = "Welcome! @%s is a DJ"
)
// WEB TEXT MESSAGES
const (
	TxtWebLoginAlready = "Too late! This party already have DJ"
	TxtWebLoginNoSuchChat = "no such chat"
)

// BOT COMMANDS
const (
	CmdLogin  = "login"
	CmdLoginFinish = "i_am_DJ_now"
	CmdLogout = "logout"
)
