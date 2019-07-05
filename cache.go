package ragnar

import (
	"github.com/go-redis/redis"
)

// UserCache wraps a UserService to provide an in-memory cache.
type UserCache struct {
	cache   map[string]*User
	service UserService
	Redis   *redis.Client
}

// NewUserCache returns a new read-through cache for service.
func NewUserCache(s UserService, r *redis.Client) *UserCache {
	return &UserCache{
		cache:   make(map[string]*User),
		service: s,
		Redis:   r,
	}
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (c *UserCache) Create(u *User) error {
	err := c.service.Create(u)
	if err != nil {
		return err
	} else if u != nil {
		c.cache[u.ID] = u
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (c *UserCache) Read(u *User) error {
	// Check the local cache first.

	if uc := c.cache[u.ID]; uc != nil {
		u = uc
		return nil
	}

	// Otherwise fetch from the underlying service.
	err := c.service.Read(u)
	if err != nil {
		return err
	} else if u != nil {
		c.cache[u.ID] = u
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (c *UserCache) ReadByEmail(u *User) error {
	// Check the local cache first.
	if uc := c.cache[u.ID]; uc != nil {
		u = uc
		return nil
	}

	// TODO: make sure the caching is workign correctly
	// Otherwise fetch from the underlying service.
	err := c.service.ReadByEmail(u)
	if err != nil {
		return err
	} else if u != nil {
		c.cache[u.ID] = u
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (c *UserCache) Update(u *User) error {
	err := c.service.Update(u)
	if err != nil {
		return err
	} else if u != nil {
		c.cache[u.ID] = u
	}

	return err
}

func (c *UserCache) Delete(u *User) error {
	err := c.service.Delete(u)
	if err != nil {
		return err
	} else if u != nil {
		delete(c.cache, u.ID)
	}

	return err
}
