package main

import (
	"github.com/alanfran/gameprofile/profile"
	"github.com/cnf/structhash"
)

// ProfileWithHash is used by the application to pass along Hashes of the last known state of a Profile.
// This is used to prevent write conflicts.
type ProfileWithHash struct {
	profile.Profile
	Hash string
}

// NewProfileWithHash creates a new ProfileWithHash given a Profile.
func NewProfileWithHash(p profile.Profile) ProfileWithHash {
	hash, _ := structhash.Hash(p, 1)
	return ProfileWithHash{
		Profile: p,
		Hash:    hash,
	}
}

// IsProfileHashValid compares a given hash to the current state of a Profile.
func (a *App) IsProfileHashValid(hash string, steamid string) bool {
	p, err := a.profiles.GetProfile(steamid)
	if err != nil {
		return false
	}

	hash2, _ := structhash.Hash(p, 1)

	return hash == hash2
}
