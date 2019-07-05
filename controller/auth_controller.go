package controller

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/RobinBaeckman/ragnar"
	"github.com/RobinBaeckman/ragnar/errors"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *ragnar.UserCache) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &ragnar.User{}
		err = mapReqJSONToUser(r, u)
		if err != nil {
			return err
		}
		if u.Email == "" ||
			u.Password == "" {
			return &errors.ErrHTTP{nil, "Missing parameters", 404}
		}
		err = c.ReadByEmail(u)
		if err != nil {
			return err
		}

		if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(u.Password)); err != nil {
			return &errors.ErrHTTP{err, "Wrong Password", 401}
		}

		buid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		uid := buid.String()

		err = c.Redis.Set(uid, u.Role, 0).Err()
		if err != nil {
			panic(err)
		}

		c := http.Cookie{
			Name:     os.Getenv("COOKIE_NAME"),
			Value:    uid,
			HttpOnly: true,
		}
		http.SetCookie(w, &c)

		b, err := json.Marshal(u)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func Logout(re *redis.Client) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		c, err := r.Cookie(os.Getenv("COOKIE_NAME"))
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
