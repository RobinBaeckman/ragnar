package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
	"github.com/RobinBaeckman/rolf/pkg/valid"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) CreateUser() func(http.ResponseWriter, *http.Request) error {
	// TODO: thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		// use thing
		u := &rolf.User{}
		err = decode(r, u)
		if err != nil {
			return err
		}

		if err := validate(u); err != nil {
			return err
		}

		u.ID = uuid.New().String()

		// TODO: make the hash in the database instead
		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		// TODO: build better Role implementation
		u.Role = "user"

		err = s.Storage.Create(u)
		if err != nil {
			return err
		}

		// TODO: Fix a better solution for this
		u.PasswordHash = nil
		u.Password = ""

		b, err := json.Marshal(u)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		url := fmt.Sprintf("%s://%s:%s/v1/users/%s", rolf.Env["PROTO"], rolf.Env["HOST"], rolf.Env["PORT"], u.ID)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusCreated)
		w.Write(b)

		return nil
	}
}

func (s *Server) ReadUser() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &rolf.User{}
		u.ID = mux.Vars(r)["id"]

		if !valid.UUID(u.ID) {
			return &rolf.Error{Code: rolf.EINVALID, Message: "Invalid UUID.", Op: rolf.Trace()}
		}

		err = s.Storage.Read(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

// TODO: make sure the password is returned too
func (s *Server) ReadAllUsers() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		us := &[]rolf.User{}
		// TODO: add ReadAll to memcache
		err = s.Storage.DB.ReadAll(us)
		if err != nil {
			return err
		}

		// Todo: create a map called mapUserToResp
		b, err := json.Marshal(us)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) UpdateUser() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &rolf.User{}
		err = decode(r, u)
		if err != nil {
			return err
		}

		u.ID = mux.Vars(r)["id"]

		// TODO: integrate with validate function
		if !valid.UUID(u.ID) {
			return &rolf.Error{Code: rolf.EINVALID, Message: fmt.Sprintf("Invalid UUID: %v", u.ID), Op: rolf.Trace()}
		}

		if err := validate(u); err != nil {
			return err
		}

		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		// TODO: build better Role implementation
		u.Role = "user"

		err = s.Storage.Update(u)
		if err != nil {
			return err
		}

		// TODO: Fix a better solution for this
		u.PasswordHash = nil
		u.Password = ""

		b, err := json.Marshal(u)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) DeleteUser() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &rolf.User{}
		u.ID = mux.Vars(r)["id"]

		if !valid.UUID(u.ID) {
			return &rolf.Error{Code: rolf.EINVALID, Message: fmt.Sprintf("Invalid UUID: %v", u.ID), Op: rolf.Trace()}
		}

		err = s.Storage.Delete(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) Login() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &rolf.User{}
		if err := decode(r, u); err != nil {
			return err
		}

		switch {
		case !valid.Email(u.Email):
			return &rolf.Error{Code: rolf.EINVALID, Message: fmt.Sprintf("Invalid email: %v", u.Email), Op: rolf.Trace()}
		case !valid.Password(u.Password):
			return &rolf.Error{Code: rolf.EINVALID, Message: fmt.Sprintf("Invalid password: %v", u.Password), Op: rolf.Trace()}
		}

		if err = s.Storage.DB.ReadByEmail(u); err != nil {
			return err
		}

		if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(u.Password)); err != nil {
			return &rolf.Error{Code: rolf.EUNAUTHORIZED, Message: "Wrong username or password.", Op: rolf.Trace(), Err: err}
		}

		et := time.Now().Add(5 * time.Minute)

		claims := &claims{
			email: u.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: et.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		c := http.Cookie{
			Name:     rolf.Env["COOKIE_NAME"],
			Value:    tokenString,
			HttpOnly: true,
			Expires:  et,
		}
		http.SetCookie(w, &c)

		u.Password = ""
		u.PasswordHash = nil

		b, err := json.Marshal(u)
		if err != nil {
			return &rolf.Error{Code: rolf.EINTERNAL, Op: rolf.Trace(), Err: err}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) Logout() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		c, err := r.Cookie(rolf.Env["COOKIE_NAME"])
		if err != nil {
			return &rolf.Error{Code: rolf.EFORBIDDEN, Message: "You are already logged out.", Op: rolf.Trace(), Err: err}
		}
		v := c.Value

		s.Storage.MemDB.Del(v)

		w.WriteHeader(http.StatusOK)

		return nil
	}
}

// TODO: make seperate functions for each validation instead.
func validate(u *rolf.User) error {
	switch {
	case !valid.Email(u.Email):
		return &rolf.Error{Code: rolf.EINVALID, Message: fmt.Sprintf("Invalid email: %v", u.Email), Op: rolf.Trace()}
	case !valid.Password(u.Password):
		return &rolf.Error{Code: rolf.EINVALID, Message: fmt.Sprintf("Invalid password: %v", u.Password), Op: rolf.Trace()}
	case !valid.FirstName(u.FirstName):
		return &rolf.Error{Code: rolf.EINVALID, Message: fmt.Sprintf("Invalid firstName: %v", u.FirstName), Op: rolf.Trace()}
	case !valid.LastName(u.LastName):
		return &rolf.Error{Code: rolf.EINVALID, Message: fmt.Sprintf("Invalid lastName: %v", u.LastName), Op: rolf.Trace()}
	}

	return nil
}
