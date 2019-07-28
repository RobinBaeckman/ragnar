package rest

import (
	"encoding/json"
	"net/http"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
)

func decode(r *http.Request, u *rolf.User) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(u)
	if err != nil {
		return &rolf.Error{Code: rolf.EINVALID, Message: "Missing parameters", Op: rolf.Trace(), Err: err}
	}

	return nil
}
