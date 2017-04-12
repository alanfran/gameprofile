package gameprofile

import "time"

// Profile stores informatin about a player.
type Profile struct {
	ID        string
	Coins     int64
	Inventory map[string]string // itemname -> settings
	Equipment map[string]string // slot -> itemname
}

// Punishment stores information about a punishment (eg: bans).
type Punishment struct {
	ID       int64
	PlayerID string
	By       string
	Type     string
	Reason   string
	Date     time.Time
	Expires  time.Time
}

// Storer defines the behavior of a Profile Store.
type Storer interface {
	GetProfile(string) (Profile, error)
	PutProfile(Profile) error

	//GetCoins(string) int64
	//PutCoins(string, int64) int64

	GetPunishments(string) (map[string]Punishment, error)
	PutPunishment(Punishment) error
	DelPunishment(int64) error
}
