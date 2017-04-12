package gameprofile

import pg "gopkg.in/pg.v4"

// PostgresStore stores a reference to the database and has methods for interacting with it.
type PostgresStore struct {
	db *pg.DB
}

// NewPostgresStore ensures the required database tables exist and returns an initialized PostgresStore.
func NewPostgresStore(db *pg.DB) *PostgresStore {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS profiles (
    id TEXT PRIMARY KEY,
    coins BIGINT,
    inventory JSONB,
		equipment JSONB
    )`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS punishments (
		id BIGSERIAL,
		player_id TEXT NOT NULL,
    by TEXT NOT NULL,
    type TEXT NOT NULL,
    reason TEXT,
    date TIMESTAMP,
    expires TIMESTAMP,
    PRIMARY KEY(id, type)
    )`)
	if err != nil {
		panic(err)
	}
	return &PostgresStore{db}
}

// GetProfile retrieves a player's profile from the database, or an empty one if not found.
func (s PostgresStore) GetProfile(playerID string) (p Profile, err error) {
	p = Profile{ID: playerID}
	err = s.db.Select(&p)
	return p, err
}

// PutProfile puts a profile into the database.
func (s PostgresStore) PutProfile(p Profile) error {
	err := s.db.Create(&p)
	if err != nil {
		_, err = s.db.Model(&p).Update()
	}
	return err
}

// GetCoins

// PutCoins

// GetPunishments returns the punishments for a player.
func (s PostgresStore) GetPunishments(steamid string) (map[string]Punishment, error) {
	var r []Punishment
	m := make(map[string]Punishment)

	err := s.db.Model(&r).Where("player_id = ?", steamid).Select()
	if err != nil {
		return m, err
	}

	for _, v := range r {
		m[v.Type] = v
	}
	return m, err
}

// PutPunishment adds a punishment to the database.
func (s PostgresStore) PutPunishment(p Punishment) error {
	err := s.db.Create(&p)

	return err
}

// DelPunishment deletes a punishment from the database.
func (s PostgresStore) DelPunishment(punishmentID int64) error {
	p := Punishment{ID: punishmentID}
	err := s.db.Delete(&p)
	return err
}
