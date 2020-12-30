# @SpotyPartyBot for Telegram
*- Dude, you should listen this one. Wait, wait, wait there will be a guitar solo*

Tired of your friend who puts his fav songs for hours at the party?

Everybody fighting for AUX cable to play the next song in a car?

Just add this bot to a chat and anyone can queue songs on a single device.
The best part? Nobody can add more than 3 songs in a row and there is a vote for skipping the song!
## How to use
1. Add `@SpotyPartyBot` to a group chat
2. Someone from the chat should log into Spotify account with `/login` command. 
This person will be a **DJ** and other people from the chat will queue songs to **DJ's** device.
   By logging in you will give a permission to bot for playing songs on your Spotify device. 
   You can take it back at any time with `/logout` command or in your [Spotify account settings](https://www.spotify.com/account/apps/) 
3. Now everybody in the chat can add their favorite songs to **DJ's** playlist. 
   Just type `@SpotyPartyBot songname` and you will see suggested song results in the same way as you search GIFs.
   *NOTE: DJ should start playing something on his Spotify device to make his device active. It will be not possible to add songs with no active device.*
4. When the party is over **DJ** can `/logout` so nobody will be able to play songs on his/her device anymore.
## Commands
Command | Description
 --- | ------
`/login` | Became a DJ. Log into Spotify and everyone in the chat will be able to queue songs on you device.
`/logout` | Stop being a DJ. Now someone else can try.
`/skip` | Start vote for skipping current song. If half of the chat will use `/skip`, current song will be skipped.
`@SpotyPartyBot ...` | Inline command to search and submit your song. Start typing song name after `@SpotyPartyBot` and you will see suggestions (after 4 characters). Nobody can submit more than 3 songs in a row.
## Roadmap
1. Show current queue state - Not supported by Spotify API 
2. `/settings` command - change allowed number of songs in a row from the same person
3. `/kick` - kick a DJ
4. show statistics: total songs, songs per person, most popular artist, etc.
5. Submit your suggestion as issue or PR