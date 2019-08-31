package redis

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
	"github.com/go-redis/redis"
)

func NewMemDB() rolf.MemDB {
	r := redis.NewClient(&redis.Options{
		Addr:     rolf.Env["REDIS_HOST"] + ":" + rolf.Env["REDIS_PORT"],
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rolf.MemDB(MemDB{r})
}

type MemDB struct {
	Redis *redis.Client
}

func (mdb MemDB) GetUser(id string) (*rolf.User, error) {
	u := &rolf.User{ID: id}
	v, err := mdb.Redis.Get(id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}
	}

	err = json.Unmarshal([]byte(v), u)
	if err != nil {
		return nil, &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
	}

	return u, nil
}

func (mdb MemDB) SetUser(k string, u *rolf.User, expiration time.Duration) (err error) {
	json, err := json.Marshal(u)
	if err != nil {
		return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
	}
	if err := mdb.Redis.Set(k, json, expiration).Err(); err != nil {
		return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
	}

	return nil
}

func (mdb MemDB) Del(v string) (err error) {
	if err := mdb.Redis.Del(v).Err(); err != nil {
		return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
	}

	return err
}
