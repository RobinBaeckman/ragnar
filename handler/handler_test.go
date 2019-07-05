package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RobinBaeckman/atla"
	"github.com/RobinBaeckman/atla/controller"
	"github.com/RobinBaeckman/atla/domain"
	"github.com/RobinBaeckman/atla/mock"
	"github.com/gorilla/mux"
)

func TestUserShow(t *testing.T) {
	// Inject our mock into our handler.
	var as mock.UserService

	// Mock our UserGet() call.
	as.UserGetFn = func(a *domain.User) error {
		if a.ID != "9463268a-0825-4ca3-b55b-75c0e2487ac1" {
			t.Fatalf("unexpected id: %s", a.ID)
		}
		return nil
	}

	// Invoke the handler.
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "localhost:3000/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	r = mux.SetURLVars(r, map[string]string{"id": "9463268a-0825-4ca3-b55b-75c0e2487ac1"})

	i := &atla.UserInteractor{
		UserService: atla.UserService(&as),
	}
	var c = controller.UserShow(i)
	err = c(w, r)
	if err != nil {
		t.Fatal(err)
	}

	// Validate mock.
	if !as.UserGetInvoked {
		t.Fatal("expected UserGet() to be invoked")
	}
}

func TestUserStore(t *testing.T) {
	// Inject our mock into our handler.
	var as mock.UserService

	// Mock our UserGet() call.
	as.UserPersistFn = func(a *domain.User) error {
		if a.Email != "email@mail.com" ||
			a.FirstName != "Tom" ||
			a.LastName != "Hanks" {
			t.Fatalf("unexpected email, firstName or lastName: %s", a.ID)
		}
		return nil
	}

	// Invoke the handler.
	w := httptest.NewRecorder()
	//pl := &struct {
	//email     string
	//password  string
	//firstName string
	//lastName  string
	//}{
	//email:     "email@mail.com",
	//password:  "secret",
	//firstName: "Tom",
	//lastName:  "Hanks",
	//}
	pl := atla.InUserStore{
		Email:     "email@mail.com",
		Password:  "secret",
		FirstName: "Tom",
		LastName:  "Hanks",
	}

	b, err := json.Marshal(&pl)

	r, err := http.NewRequest("POST", "localhost:3000/v1/users", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	i := &atla.UserInteractor{
		UserService: atla.UserService(&as),
	}
	var c = controller.UserStore(i)
	err = c(w, r)
	if err != nil {
		t.Fatal(err)
	}

	// Validate mock.
	if !as.UserPersistInvoked {
		t.Fatal("expected UserPersist() to be invoked")
	}
}
