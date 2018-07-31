package controller

import (
	"encoding/json"
	"net/http"

	"github.com/RobinBaeckman/ragnar"
	"github.com/RobinBaeckman/ragnar/errors"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *ragnar.UserCache) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		in, err := decode(r)
		if err != nil {
			return err
		}
		if in.Email == "" ||
			in.Password == "" {
			return &errors.ErrHTTP{nil, "Missing parameters", 404}
		}
		u, err := c.GetByEmail(in.Email)

		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(in.Password)); err != nil {
			return &errors.ErrHTTP{err, "Wrong Password", 401}
		}

		uid := uuid.NewV4().String()
		err = c.Redis.Set(uid, u.Role, 0).Err()
		if err != nil {
			panic(err)
		}

		c := http.Cookie{
			Name:     viper.GetString("session.cookie_name"),
			Value:    uid,
			HttpOnly: true,
		}
		http.SetCookie(w, &c)

		vm := &UserViewModel{}
		vm.Map(u)

		j, err := json.Marshal(vm)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)

		return nil
	}
}

func Logout(re *redis.Client) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		c, err := r.Cookie(viper.GetString("session.cookie_name"))
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		v := c.Value

		// Check if user is authenticated
		re.Del(v)

		w.WriteHeader(http.StatusOK)

		return nil
	}
}
