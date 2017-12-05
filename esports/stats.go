package esports

type GameStats struct {
	GameID       int64
	PlatformID   string
	GameCreation int64 // TODO: factor out apiclient/unix_millis and use it here.
	GameDuration int64
	QueueID      int
	MapID        int
	SeasonID     int
	GameVersion  string
	GameMode     string
	GameType     string
}
