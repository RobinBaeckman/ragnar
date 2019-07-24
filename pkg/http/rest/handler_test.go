package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
	"time"

	uuid "github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/RobinBaeckman/ragnar/pkg/http/rest"
	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
	"github.com/RobinBaeckman/ragnar/pkg/valid"
)

var s *rest.Server
var us *[]ragnar.User
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	dbenv := os.Getenv("MYSQL_DB")
	// TODO: change from static to env test db
	os.Setenv("MYSQL_DB", "ragnar_db_test")
	err := rest.ParseEnv()
	if err != nil {
		panic(err)
	}

	s, err = rest.NewServer()
	if err != nil {
		panic(err)
	}

	s.Storage.DB.CleanupTables()

	us, err = createRandomUsers(2)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Setenv("MYSQL_DB", dbenv)
	s.Storage.DB.CleanupTables()
	s.Storage.DB.Close()

	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	reqB := &ragnar.User{}
	generateUserData(reqB)

	b, err := json.Marshal(reqB)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: change from static to env hostname
	r, err := http.NewRequest("POST", "localhost:3000/v1/users", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	h := s.LogAndError(s.CreateUser())

	w := httptest.NewRecorder()
	err = h(w, r)
	if err != nil {
		t.Fatal(err)
	}

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("Wrong status code: got %v want %v", status, http.StatusCreated)
	}

	resB := &ragnar.User{}

	d := json.NewDecoder(w.Body)
	err = d.Decode(resB)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: make this a function instead
	switch {
	case !valid.UUID(resB.ID):
		t.Errorf("Invalid uuid: %v", resB.ID)
	case resB.Email != reqB.Email:
		t.Errorf("Invalid parameters: %v or %v", resB.Email, reqB.Email)
	case resB.Password != "":
		t.Errorf("Password should be empty, for security reasons")
	case resB.FirstName != reqB.FirstName:
		t.Errorf("Invalid parameters: %v or %v", resB.FirstName, reqB.FirstName)
	case resB.LastName != reqB.LastName:
		t.Errorf("Invalid parameters: %v or %v", resB.FirstName, reqB.FirstName)
	}

	*us = append(*us, *resB)
}

func TestReadUser(t *testing.T) {
	r, err := http.NewRequest("GET", "localhost:3000/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	r = mux.SetURLVars(r, map[string]string{
		"id": (*us)[0].ID,
	})

	h := s.LogAndError(s.ReadUser())

	w := httptest.NewRecorder()
	err = h(w, r)
	if err != nil {
		t.Fatal(err)
	}

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
	}

	resB := &ragnar.User{}

	d := json.NewDecoder(w.Body)
	err = d.Decode(resB)
	if err != nil {
		t.Fatal(err)
	}

	switch {
	case resB.ID != (*us)[0].ID:
		t.Errorf("Invalid parameters: %v or %v", resB.ID, (*us)[0].ID)
	case resB.Email != (*us)[0].Email:
		t.Errorf("Invalid parameters: %v or %v", resB.Email, (*us)[0].Email)
	case resB.Password != "":
		t.Errorf("Password should be empty for security reasons")
	case resB.FirstName != (*us)[0].FirstName:
		t.Errorf("Invalid parameters: %v or %v", resB.FirstName, (*us)[0].FirstName)
	case resB.LastName != (*us)[0].LastName:
		t.Errorf("Invalid parameters: %v or %v", resB.FirstName, (*us)[0].FirstName)
	}
}

func TestReadAllUsers(t *testing.T) {
	r, err := http.NewRequest("GET", "localhost:3000/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// add all middleware that usually is there to get a more genuine test
	h := s.LogAndError(s.ReadAllUsers())

	w := httptest.NewRecorder()
	err = h(w, r)
	if err != nil {
		t.Fatal(err)
	}

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
	}

	resBs := &[]ragnar.User{}

	d := json.NewDecoder(w.Body)
	err = d.Decode(resBs)
	if err != nil {
		t.Fatal(err)
	}
	sort.SliceStable(*resBs, func(i, j int) bool {
		return (*resBs)[i].Email < (*resBs)[j].Email
	})

	sort.SliceStable(*us, func(i, j int) bool {
		return (*us)[i].Email < (*us)[j].Email
	})

	for i, resB := range *resBs {
		switch {
		case resB.ID != (*us)[i].ID:
			t.Errorf("Invalid parameters: %v or %v", resB.ID, (*us)[i].ID)
		case resB.Email != (*us)[i].Email:
			t.Errorf("Invalid parameters: %v or %v", resB.Email, (*us)[i].Email)
		case resB.Password != "":
			t.Errorf("Password should be empty, for security reasons")
		case resB.FirstName != (*us)[i].FirstName:
			t.Errorf("Invalid parameters: %v or %v", resB.FirstName, (*us)[i].FirstName)
		case resB.LastName != (*us)[i].LastName:
			t.Errorf("Invalid parameters: %v or %v", resB.FirstName, (*us)[i].FirstName)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	reqB := &ragnar.User{}
	generateUserData(reqB)

	b, err := json.Marshal(reqB)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest("PUT", "localhost:3000/v1/users", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	r = mux.SetURLVars(r, map[string]string{
		"id": (*us)[1].ID,
	})

	h := s.LogAndError(s.UpdateUser())

	w := httptest.NewRecorder()
	err = h(w, r)
	if err != nil {
		t.Fatal(err)
	}

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", status, http.StatusOK)
	}

	resB := &ragnar.User{}

	d := json.NewDecoder(w.Body)
	err = d.Decode(resB)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resB)

	switch {
	case !valid.UUID(resB.ID):
		t.Errorf("Invalid uuid: %v", resB.ID)
	case resB.Email != reqB.Email:
		t.Errorf("Invalid parameters: %v or %v", resB.Email, reqB.Email)
	case resB.Password != "":
		t.Errorf("Password should be empty, for security reasons")
	case resB.FirstName != reqB.FirstName:
		t.Errorf("Invalid parameters: %v or %v", resB.FirstName, reqB.FirstName)
	case resB.LastName != reqB.LastName:
		t.Errorf("Invalid parameters: %v or %v", resB.FirstName, reqB.FirstName)
	}

	*us = append(*us, *resB)
}

func BenchmarkCreateUser(b *testing.B) {
	reqB := &ragnar.User{}
	generateUserData(reqB)

	bt, err := json.Marshal(reqB)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	// TODO: change from static to env hostname
	r, err := http.NewRequest("POST", "localhost:3000/v1/users", bytes.NewBuffer(bt))
	if err != nil {
		b.Fatal(err)
	}

	h := s.LogAndError(s.CreateUser())

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		err = h(w, r)
	}
	if err != nil {
		b.Fatal(err)
	}
}

func generateUserData(u *ragnar.User) {
	u.ID = uuid.New().String()
	u.Email = fmt.Sprintf("%s@mail.com", randSeq(12))
	u.Password = "secret"
	u.PasswordHash = []byte("$2a$10$flStfMMZw4TsuJh3OdJhYeDBCibDlTNNm.yVMya4RgMcc7bF0/2nq")
	u.FirstName = "Rolf"
	u.LastName = "Baeckman"
	u.Role = "user"
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// i is the number of random users that is being created
func createRandomUsers(i int) (*[]ragnar.User, error) {
	us := []ragnar.User{}
	for y := 0; y < i; y++ {
		u := &ragnar.User{}
		generateUserData(u)
		err := s.Storage.DB.Create(u)
		if err != nil {
			return &us, err
		}
		us = append(us, *u)
	}

	return &us, nil
}
