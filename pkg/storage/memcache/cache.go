package memcache

import (
	"github.com/RobinBaeckman/rolf/pkg/rolf"
)

// UserStorage wraps a UserService to provide an in-memory cache.
type Storage struct {
	DB    rolf.DB
	MemDB rolf.MemDB
}

// NewUserCache returns a new read-through cache for service.
// TODO: change all fields to correct visuality
func NewStorage(db rolf.DB, mdb rolf.MemDB) *Storage {
	return &Storage{
		DB:    db,
		MemDB: mdb,
	}
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) Create(u *rolf.User) error {
	err := s.DB.Create(u)
	if err != nil {
		return err
	} else if u != nil {
		if err := s.MemDB.SetUser(u.ID, u, 0); err != nil {
			return err
		}
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) Read(u *rolf.User) error {
	// Check the local cache first.
	uc, err := s.MemDB.GetUser(u.ID)
	if err != nil {
		return err
	}
	if uc != nil {
		*u = *uc
		return nil
	}

	// Otherwise fetch from the underlying service.
	err = s.DB.Read(u)
	if err != nil {
		return err
	} else if u != nil {
		if err := s.MemDB.SetUser(u.ID, u, 0); err != nil {
			return err
		}
	}

	return nil
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) ReadAny(u *rolf.User) error {
	// Check the local cache first.
	uc, err := s.MemDB.GetUser(u.ID)
	if err != nil {
		return err
	}
	if uc != nil {
		*u = *uc
		return nil
	}

	// Otherwise fetch from the underlying service.
	err = s.DB.ReadAny(u)
	if err != nil {
		return err
	} else if u != nil {
		if err := s.MemDB.SetUser(u.ID, u, 0); err != nil {
			return err
		}
	}

	return nil
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) ReadByEmail(u *rolf.User) error {
	// Check the local cache first.
	uc, err := s.MemDB.GetUser(u.ID)
	if err != nil {
		return err
	}
	if uc != nil {
		*u = *uc
		return nil
	}

	// Otherwise fetch from the underlying service.
	err = s.DB.ReadByEmail(u)
	if err != nil {
		return err
	} else if u != nil {
		if err := s.MemDB.SetUser(u.ID, u, 0); err != nil {
			return err
		}
	}

	return nil
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) Update(u *rolf.User) error {
	err := s.DB.Update(u)
	if err != nil {
		return err
	} else if u != nil {
		if err := s.MemDB.SetUser(u.ID, u, 0); err != nil {
			return err
		}
	}

	return err
}

func (s *Storage) Delete(u *rolf.User) error {
	err := s.DB.Delete(u)
	if err != nil {
		return err
	}

	if err := s.MemDB.Del(u.ID); err != nil {
		return err
	}

	return err
}
