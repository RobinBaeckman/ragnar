package redis

import (
	"net/http"
	"time"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
	"github.com/go-redis/redis"
)

func NewMemDB() MemDB {
	r := redis.NewClient(&redis.Options{
		Addr:     rolf.Env["REDIS_HOST"] + ":" + rolf.Env["REDIS_PORT"],
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return MemDB{r}
}

type MemDB struct {
	Redis *redis.Client
}

func (mdb MemDB) Get(k string) (string, error) {
	v, err := mdb.Redis.Get(k).Result()
	if err != nil {
		return v, &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
	}

	return v, err
}

func (mdb MemDB) Set(k string, v interface{}, expiration time.Duration) error {
	if err := mdb.Redis.Set(k, v, 0).Err(); err != nil {
		return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
	}

	return nil
}

func (mdb MemDB) Del(v string) error {
	if err := mdb.Redis.Del(v).Err(); err != nil {
		return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
	}

	return nil
}
