package rest

import (
	"encoding/json"
	"net/http"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
)

func mapReqJSONToUser(r *http.Request, u *ragnar.User) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(u)
	if err != nil {
		return &ragnar.Error{Code: ragnar.EINVALID, Message: "Missing parameters", Op: ragnar.Trace(), Err: err}
	}

	return nil
}
