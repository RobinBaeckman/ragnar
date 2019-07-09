package memcache

import (
	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
	"github.com/go-redis/redis"
)

// UserStorage wraps a UserService to provide an in-memory cache.
type UserStorage struct {
	Memcache map[string]*ragnar.User
	DB       ragnar.UserDB
	Redis    *redis.Client
}

// NewUserCache returns a new read-through cache for service.
func NewUserStorage(db ragnar.UserDB, r *redis.Client) *UserStorage {
	return &UserStorage{
		Memcache: make(map[string]*ragnar.User),
		DB:       db,
		Redis:    r,
	}
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *UserStorage) Create(u *ragnar.User) error {
	err := s.DB.Create(u)
	if err != nil {
		return err
	} else if u != nil {
		s.Memcache[u.ID] = u
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *UserStorage) Read(u *ragnar.User) error {
	// Check the local cache first.

	if uc := s.Memcache[u.ID]; uc != nil {
		u = uc
		return nil
	}

	// Otherwise fetch from the underlying service.
	err := s.DB.Read(u)
	if err != nil {
		return err
	} else if u != nil {
		s.Memcache[u.ID] = u
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *UserStorage) ReadByEmail(u *ragnar.User) error {
	// Check the local cache first.
	if uc := s.Memcache[u.ID]; uc != nil {
		u = uc
		return nil
	}

	// TODO: make sure the caching is workign correctly
	// Otherwise fetch from the underlying service.
	err := s.DB.ReadByEmail(u)
	if err != nil {
		return err
	} else if u != nil {
		s.Memcache[u.ID] = u
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (s *UserStorage) Update(u *ragnar.User) error {
	err := s.DB.Update(u)
	if err != nil {
		return err
	} else if u != nil {
		s.Memcache[u.ID] = u
	}

	return err
}

func (s *UserStorage) Delete(u *ragnar.User) error {
	err := s.DB.Delete(u)
	if err != nil {
		return err
	} else if u != nil {
		delete(s.Memcache, u.ID)
	}

	return err
}
