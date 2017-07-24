package profile

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

type BoltStore struct {
	db *bolt.DB
}

func NewBoltStore(path string) (*BoltStore, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Ensure the profiles and punishments buckets exist.
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("profiles"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("punishments"))
		if err != nil {
			return err
		}

		return nil
	})

	return &BoltStore{db: db}, nil
}

// PutProfile stores the JSON representation of a Profile in the database with its ID as the key.
func (s *BoltStore) PutProfile(p Profile) error {
	j, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("profiles"))
		err := b.Put([]byte(p.ID), j)
		return err
	})
}

// GetProfile retrieves a Profile from the database with a key that matches the steamid.
func (s *BoltStore) GetProfile(steamid string) (Profile, error) {
	var p Profile

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("profiles"))
		v := b.Get([]byte(steamid))
		if v == nil {
			return errors.New("Profile not found.")
		}
		return json.Unmarshal(v, &p)
	})

	return p, err
}

//GetCoins(string) int64
//PutCoins(string, int64) int64

func (s *BoltStore) GetPunishments(steamid string) (map[string]Punishment, error) {
	var ps map[string]Punishment

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("punishments"))
		v := b.Get([]byte(steamid))
		if v == nil {
			return errors.New("No punishments found for SteamID " + steamid)
		}
		return json.Unmarshal(v, &ps)
	})

	return ps, err
}

func (s *BoltStore) PutPunishment(p Punishment) error {
	var ps map[string]Punishment

	ps, err := s.GetPunishments(p.PlayerID)
	if err != nil {
		ps = map[string]Punishment{}
	}

	ps[p.Type] = p

	j, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("punishments"))
		return b.Put([]byte(p.PlayerID), j)
	})

	return err
}

func (s *BoltStore) DelPunishment(pid int64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("punishments"))

		return b.ForEach(func(k, v []byte) error {
			var ps map[string]Punishment
			err := json.Unmarshal(v, &ps)
			if err != nil {
				return err
			}
			for _, v := range ps {
				if v.ID == pid {
					delete(ps, v.Type)
					j, err := json.Marshal(ps)
					if err != nil {
						return err
					}
					err = b.Put([]byte(v.PlayerID), j)
					return err
				}
			}
			return err
		})
	})
}
