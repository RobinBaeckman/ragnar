package memcache

import (
	"sync"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
)

func (s *Storage) memGet(key string) *rolf.User {
	s.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer s.mux.Unlock()
	return s.Memcache[key]
}

func (s *Storage) memSet(u *rolf.User) {
	s.mux.Lock()
	s.Memcache[u.ID] = u
	s.mux.Unlock()
}

func (s *Storage) memDel(key string) {
	s.mux.Lock()
	delete(s.Memcache, key)
	s.mux.Unlock()
}

// UserStorage wraps a UserService to provide an in-memory cache.
type Storage struct {
	Memcache map[string]*rolf.User
	DB       rolf.DB
	MemDB    rolf.MemDB
	mux      sync.Mutex
}

// NewUserCache returns a new read-through cache for service.
// TODO: change all fields to correct visuality
func NewStorage(db rolf.DB, mdb rolf.MemDB) *Storage {
	return &Storage{
		Memcache: make(map[string]*rolf.User),
		DB:       db,
		MemDB:    mdb,
	}
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) Create(u *rolf.User) error {
	err := s.DB.Create(u)
	if err != nil {
		return err
	} else if u != nil {
		s.memSet(u)
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) Read(u *rolf.User) error {
	// Check the local cache first.

	if uc := s.memGet(u.ID); uc != nil {
		u = uc
		return nil
	}

	// Otherwise fetch from the underlying service.
	err := s.DB.Read(u)
	if err != nil {
		return err
	} else if u != nil {
		s.memSet(u)
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) ReadByEmail(u *rolf.User) error {
	// Check the local cache first.
	if uc := s.memGet(u.ID); uc != nil {
		u = uc
		return nil
	}

	// TODO: make sure the caching is workign correctly
	// Otherwise fetch from the underlying service.
	err := s.DB.ReadByEmail(u)
	if err != nil {
		return err
	} else if u != nil {
		s.memSet(u)
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *Storage) Update(u *rolf.User) error {
	err := s.DB.Update(u)
	if err != nil {
		return err
	} else if u != nil {
		s.memSet(u)
	}

	return err
}

func (s *Storage) Delete(u *rolf.User) error {
	err := s.DB.Delete(u)
	if err != nil {
		return err
	} else if u != nil {
		s.memDel(u.ID)
	}

	return err
}
