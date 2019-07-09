package rest

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) CreateUser() func(http.ResponseWriter, *http.Request) error {
	// TODO: thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		// use thing
		u := &ragnar.User{}
		err = mapReqJSONToUser(r, u)
		if err != nil {
			return err
		}

		// Validation
		if u.Email == "" ||
			u.Password == "" ||
			u.LastName == "" ||
			u.FirstName == "" {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: "Missing parameters", Op: ragnar.Trace(), Err: err}
		}
		buid, err := uuid.NewV4()
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}
		u.ID = buid.String()
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

		if u.ID == "" {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: "Missing parameters", Op: ragnar.Trace(), Err: err}
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
		err = mapReqJSONToUser(r, u)
		if err != nil {
			return err
		}

		u.ID = mux.Vars(r)["id"]

		// Validation
		// TODO: improve validation
		// TODO: u.ID should be a valid uuid
		if u.ID == "" ||
			u.Email == "" ||
			u.Password == "" ||
			u.LastName == "" ||
			u.FirstName == "" {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: "Missing parameters", Op: ragnar.Trace(), Err: err}
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

		// Validation
		// TODO: improve validation
		// TODO: u.ID should be a valid uuid
		if u.ID == "" {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: "Missing parameters", Op: ragnar.Trace(), Err: err}
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
		err = mapReqJSONToUser(r, u)
		if err != nil {
			return err
		}
		if u.Email == "" ||
			u.Password == "" {
			return &ragnar.Error{Code: ragnar.EINVALID, Message: "Missing parameters", Op: ragnar.Trace(), Err: err}
		}
		err = s.userStorage.DB.ReadByEmail(u)
		if err != nil {
			return err
		}

		if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(u.Password)); err != nil {
			return &ragnar.Error{Code: ragnar.EUNAUTHORIZED, Message: "Wrong username or password.", Op: ragnar.Trace(), Err: err}
		}

		buid, err := uuid.NewV4()
		if err != nil {
			return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
		}
		u.ID = buid.String()

		err = s.userStorage.Redis.Set(u.ID, u.Email, 0).Err()
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
