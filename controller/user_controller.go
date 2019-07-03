package controller

import (
	"encoding/json"
	"net/http"

	"github.com/RobinBaeckman/ragnar"
	"github.com/RobinBaeckman/ragnar/errors"
)

type (
	UserViewModel struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}
)

func CreateUser(c *ragnar.UserCache) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		in, err := decode(r)
		if err != nil {
			return err
		}

		if in.Email == "" ||
			in.Password == "" ||
			in.LastName == "" ||
			in.FirstName == "" {
			return &errors.ErrHTTP{nil, "Missing parameters", 404}
		}

		u, err := mapUserStore(&in)
		if err != nil {
			return err
		}

		err = c.Create(&u)
		if err != nil {
			return err
		}

		vm := &UserViewModel{}
		vm.Map(&u)

		j, err := json.Marshal(vm)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)

		return nil
	}
}

func ReadUser(c *ragnar.UserCache) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		u, err := mapUserShow(r)
		if err != nil {
			return err
		}

		if u.ID == "" {
			return &errors.ErrHTTP{nil, "Missing parameters", 404}
		}

		err = c.Read(&u)
		if err != nil {
			return err
		}
		vm := &UserViewModel{}
		vm.Map(&u)

		j, err := json.Marshal(vm)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)

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
		vms := []UserViewModel{}
		for _, u := range *us {
			vm := &UserViewModel{}
			vm.Map(&u)
			vms = append(vms, *vm)
		}

		j, err := json.Marshal(vms)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)

		return nil

	}

}

func (vm *UserViewModel) Map(u *ragnar.User) {
	*vm = UserViewModel{
		ID:        u.ID,
		Email:     u.Email,
		Password:  "*****",
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}
