package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RobinBaeckman/ragnar"
	"github.com/RobinBaeckman/ragnar/errors"
	"github.com/gorilla/mux"
)

func CreateUser(c *ragnar.UserCache) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
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
			return &errors.ErrHTTP{nil, "Missing parameters", 404}
		}

		err = c.Create(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func ReadUser(c *ragnar.UserCache) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &ragnar.User{}
		u.ID = mux.Vars(r)["id"]

		if u.ID == "" {
			return &errors.ErrHTTP{nil, "Missing parameters", 404}
		}

		err = c.Read(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func ReadAllUsers(s ragnar.UserService) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		us := &[]ragnar.User{}
		err = s.ReadAll(us)
		if err != nil {
			return err
		}

		// Todo: create a map called mapUserToResp
		b, err := json.Marshal(us)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func UpdateUser(c *ragnar.UserCache) func(http.ResponseWriter, *http.Request) error {
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
			return &errors.ErrHTTP{nil, "Missing parameters", 404}
		}

		err = c.Update(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func DeleteUser(c *ragnar.UserCache) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u := &ragnar.User{}
		u.ID = mux.Vars(r)["id"]

		// Validation
		// TODO: improve validation
		// TODO: u.ID should be a valid uuid
		if u.ID == "" {
			return &errors.ErrHTTP{nil, "Missing parameters", 404}
		}

		err = c.Delete(u)
		if err != nil {
			return err
		}

		b, err := json.Marshal(u)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}
