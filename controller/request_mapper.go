package controller

import (
	"encoding/json"
	"net/http"

	"github.com/RobinBaeckman/ragnar"
)

func mapReqJSONToUser(r *http.Request, u *ragnar.User) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(u)
	if err != nil {
		return err
	}

	return nil
}
