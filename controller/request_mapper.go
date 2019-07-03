package controller

import (
	"encoding/json"
	"net/http"

	"github.com/RobinBaeckman/ragnar"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func decode(r *http.Request) (ragnar.In, error) {
	in := ragnar.In{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&in)
	if err != nil {
		return in, err
	}

	return in, nil
}

func mapUserStore(in *ragnar.In) (ragnar.User, error) {
	u := ragnar.User{}
	buid, err := uuid.NewV4()
	if err != nil {
		return u, err
	}
	uid := buid.String()

	pHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return u, err
	}

	u.ID = uid
	u.Email = in.Email
	u.Password = pHash
	u.FirstName = in.FirstName
	u.LastName = in.LastName
	u.Role = in.Role

	return u, nil
}

func mapUserShow(r *http.Request) (ragnar.User, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	u := ragnar.User{
		ID: id,
	}

	return u, nil
}
