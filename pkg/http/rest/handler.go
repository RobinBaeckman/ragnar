package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
	"github.com/RobinBaeckman/ragnar/pkg/valid"
	uuid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) CreateUser() func(http.ResponseWriter, *http.Request) error {
	// TODO: thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		// use thing
		u := &ragnar.User{}
		err = decode(r, u)
		if err != nil {
			return err
		}

		var msg string
		switch {
		case !valid.IsEmail(u.Email):
			msg = fmt.Sprintf("Invalid parameter: %v", u.Email)
		case !valid.IsPassword(u.Password):
			msg = fmt.Sprintf("Invalid parameter: %v", u.Password)
		case !valid.IsFirstName(u.FirstName):
			msg = fmt.Sprintf("Invalid parameter: %v", u.FirstName)
		case !valid.IsLastName(u.LastName):
			msg = fmt.Sprintf("Invalid parameter: %v", u.LastName)
		}
		if msg != "" {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: msg, Op: ragnar.Trace()}
		}

		u.ID = uuid.New().String()

		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}

		// TODO: build better Role implementation
		u.Role = "user"

		err = s.userStorage.Create(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) ReadUser() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &ragnar.User{}
		u.ID = mux.Vars(r)["id"]

		if !valid.IsUUID(u.ID) {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: "Invalid UUID.", Op: ragnar.Trace()}
		}

		err = s.userStorage.Read(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) ReadAllUsers() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		us := &[]ragnar.User{}
		// TODO: add ReadAll to memcache
		err = s.userStorage.DB.ReadAll(us)
		if err != nil {
			return err
		}

		// Todo: create a map called mapUserToResp
		b, err := json.Marshal(us)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) UpdateUser() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &ragnar.User{}
		err = decode(r, u)
		if err != nil {
			return err
		}

		u.ID = mux.Vars(r)["id"]

		var msg string
		switch {
		case !valid.IsUUID(u.ID):
			msg = fmt.Sprintf("Invalid parameter: %v", u.ID)
		case !valid.IsEmail(u.Email):
			msg = fmt.Sprintf("Invalid parameter: %v", u.Email)
		case !valid.IsPassword(u.Password):
			msg = fmt.Sprintf("Invalid parameter: %v", u.Password)
		case !valid.IsFirstName(u.FirstName):
			msg = fmt.Sprintf("Invalid parameter: %v", u.FirstName)
		case !valid.IsLastName(u.LastName):
			msg = fmt.Sprintf("Invalid parameter: %v", u.LastName)
		}
		if msg != "" {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: msg, Op: ragnar.Trace()}
		}

		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}

		// TODO: build better Role implementation
		u.Role = "user"

		err = s.userStorage.Update(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) DeleteUser() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &ragnar.User{}
		u.ID = mux.Vars(r)["id"]

		if !valid.IsUUID(u.ID) {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: "Invalid UUID.", Op: ragnar.Trace()}
		}

		err = s.userStorage.Delete(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) Login() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &ragnar.User{}
		if err := decode(r, u); err != nil {
			return err
		}

		var msg string
		switch {
		case !valid.IsEmail(u.Email):
			msg = fmt.Sprintf("Invalid parameter: %v", u.Email)
		case !valid.IsPassword(u.Password):
			msg = fmt.Sprintf("Invalid parameter: %v", u.Password)
		}
		if msg != "" {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: msg, Op: ragnar.Trace()}
		}

		if err = s.userStorage.DB.ReadByEmail(u); err != nil {
			return err
		}

		if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(u.Password)); err != nil {
			return &ragnar.Error{Code: ragnar.EUNAUTHORIZED, Message: "Wrong username or password.", Op: ragnar.Trace(), Err: err}
		}

		uid := uuid.New().String()

		err = s.userStorage.Redis.Set(uid, u.Email, 0).Err()
		if err != nil {
			return &ragnar.Error{Code: ragnar.EUNAUTHORIZED, Op: ragnar.Trace(), Err: err}
		}

		c := http.Cookie{
			Name:     os.Getenv("COOKIE_NAME"),
			Value:    u.ID,
			HttpOnly: true,
		}
		http.SetCookie(w, &c)

		b, err := json.Marshal(u)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) Logout() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		c, err := r.Cookie(os.Getenv("COOKIE_NAME"))
		if err != nil {
			return &ragnar.Error{Code: ragnar.EFORBIDDEN, Message: "You are already logged out.", Op: ragnar.Trace(), Err: err}
		}
		v := c.Value

		s.userStorage.Redis.Del(v)

		w.WriteHeader(http.StatusOK)

		return nil
	}
}
