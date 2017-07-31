package profile

import "errors"

// MockStore provides a simple in-memory store for use in unit tests.
type MockStore struct {
	profiles          map[string]Profile
	punishments       map[int64]Punishment
	punishmentsSerial int64
}

// NewMockStore returns an initialized MockStore.
func NewMockStore() *MockStore {
	return &MockStore{
		profiles:    map[string]Profile{},
		punishments: map[int64]Punishment{},
	}
}

// GetProfile returns a profile with the matching ID.
func (s *MockStore) GetProfile(id string) (p Profile, err error) {
	p, ok := s.profiles[id]
	if !ok {
		return p, errors.New("Profile not found.")
	}
	return p, nil
}

// PutProfile stores a profile.
func (s *MockStore) PutProfile(p Profile) error {
	if p.ID == "" {
		return errors.New("Error putting profile: no ID provided.")
	}

	s.profiles[p.ID] = p
	return nil
}

//GetCoins(string) int64
//PutCoins(string, int64) int64

// GetPunishments return a user's punishments.
func (s *MockStore) GetPunishments(pid string) (ps map[string]Punishment, err error) {
	ps = map[string]Punishment{}

	for k := range s.punishments {
		if s.punishments[k].PlayerID == pid {
			ps[s.punishments[k].Type] = s.punishments[k]
		}
	}

	if len(ps) == 0 {
		err = errors.New("No punishments found.")
	}

	return ps, err
}

// PutPunishment stores a punishment.
func (s *MockStore) PutPunishment(p Punishment) error {
	if p.PlayerID == "" || p.By == "" || p.Type == "" {
		return errors.New("PlayerID, By, and Reason fields are required.")
	}

	if p.ID == 0 {
		p.ID = s.punishmentsSerial
		s.punishmentsSerial++
	}
	s.punishments[p.ID] = p

	return nil
}

// DelPunishment removes a punishment from the store.
func (s *MockStore) DelPunishment(id int64) error {
	_, ok := s.punishments[id]
	if !ok {
		return errors.New("Error deleting punishment: ID not found")
	}

	delete(s.punishments, id)
	return nil
}
