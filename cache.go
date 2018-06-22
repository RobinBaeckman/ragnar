package ragnar

import "github.com/go-redis/redis"

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
func (c *UserCache) Get(u *User) error {
	// Check the local cache first.

	if uc := c.cache[u.ID]; uc != nil {
		u = uc
		return nil
	}

	// Otherwise fetch from the underlying service.
	err := c.service.Get(u)
	if err != nil {
		return err
	} else if u != nil {
		c.cache[u.ID] = u
	}

	return err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (c *UserCache) GetByEmail(e string) (*User, error) {
	// Check the local cache first.
	if uc := c.cache[e]; uc != nil {
		return uc, nil
	}

	// Otherwise fetch from the underlying service.
	u, err := c.service.GetByEmail(e)
	if err != nil {
		return u, err
	} else if u != nil {
		c.cache[u.Email] = u
	}

	return u, err
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (c *UserCache) Store(u *User) error {
	err := c.service.Store(u)
	if err != nil {
		return err
	} else if u != nil {
		c.cache[u.ID] = u
	}

	return err
}
