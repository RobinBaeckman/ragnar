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

	"github.com/RobinBaeckman/rolf/pkg/http/rest"
	"github.com/RobinBaeckman/rolf/pkg/rolf"
	"github.com/RobinBaeckman/rolf/pkg/valid"
)

var s *rest.Server
var us *[]rolf.User
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	dbenv := os.Getenv("MYSQL_DB")
	// TODO: change from static to env test db
	os.Setenv("MYSQL_DB", "rolf_db_test")
	err := rest.ParseEnv()
	if err != nil {
		panic(err)
	}

	s, err = rest.NewServer()
	if err != nil {
		panic(err)
	}

	s.Storage.DB.CleanupTables()

	code := m.Run()

	os.Setenv("MYSQL_DB", dbenv)
	s.Storage.DB.CleanupTables()
	s.Storage.DB.Close()

	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	tests := map[string]struct {
		ID         string `json="id,omitempty"`
		Email      string `json="email,omitempty"`
		Password   string `json="password,omitempty"`
		FirstName  string `json="firstName,omitempty"`
		LastName   string `json="lastName,omitempty"`
		wantedCode int    `json="-"`
	}{
		"simple":                 {Email: "user1@mail.com", Password: "secret", FirstName: "Rolf", LastName: "Baeckman", wantedCode: 201},
		"email misformatted":     {Email: "user2mail.com", Password: "secret", FirstName: "Rolf", LastName: "Baeckman", wantedCode: 422},
		"firstname misformatted": {Email: "user3@mail.com", Password: "secret", FirstName: "Rolf#$", LastName: "Baeckman", wantedCode: 422},
		"lastname misformatted":  {Email: "user4@mail.com", Password: "secret", FirstName: "Rolf", LastName: "Baeckman$%", wantedCode: 422},
		"password to short":      {Email: "user5@mail.com", Password: "sec", FirstName: "Rolf", LastName: "Baeckman", wantedCode: 422},
		"email missing":          {Password: "secret", FirstName: "Rolf", LastName: "Baeckman", wantedCode: 422},
		"password missing":       {Email: "user6@mail.com", FirstName: "Rolf", LastName: "Baeckman", wantedCode: 422},
		"firstname missing":      {Email: "user7@mail.com", Password: "secret", LastName: "Baeckman", wantedCode: 422},
		"lastname missing":       {Email: "user8@mail.com", Password: "secret", FirstName: "Rolf", wantedCode: 422},
	}

	for name, reqB := range tests {
		t.Run(name, func(t *testing.T) {
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

			switch code := w.Code; {
			case code != reqB.wantedCode:
				t.Fatalf("Wrong status code: got %v want %v", code, reqB.wantedCode)
			case code != 201:
				return
			}

			resB := &struct {
				ID        string `json="id"`
				Email     string `json="email"`
				Password  string `json="password"`
				FirstName string `json="firstName"`
				LastName  string `json="lastName"`
				Role      string `json="role"`
			}{}

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
				t.Errorf("Invalid role: %v or %v", resB.FirstName, reqB.FirstName)
			}
		})
	}

	s.Storage.DB.CleanupTables()
}

func TestReadUser(t *testing.T) {
	id := "90cf445a-b18e-4f29-9139-57fab62aedbc"
	tests := map[string]struct {
		id         string
		wantedCode int `json="-"`
	}{
		"simple":         {id: id, wantedCode: 200},
		"id misformated": {id: "$%)#(%'%')", wantedCode: 422},
		"id empty":       {id: "", wantedCode: 422},
	}

	ph := []byte("$2a$10$flStfMMZw4TsuJh3OdJhYeDBCibDlTNNm.yVMya4RgMcc7bF0/2nq")
	u := &rolf.User{ID: id, Email: "readUser@mail.com", Password: "secret", PasswordHash: ph, FirstName: "Rolf", LastName: "Baeckman", Role: "user"}

	err := s.Storage.DB.Create(u)
	if err != nil {
		t.Fatal(err)
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "localhost:3000/v1/users", nil)
			if err != nil {
				t.Fatal(err)
			}

			r = mux.SetURLVars(r, map[string]string{
				"id": tc.id,
			})

			h := s.LogAndError(s.ReadUser())

			w := httptest.NewRecorder()
			err = h(w, r)

			switch code := w.Code; {
			case code != tc.wantedCode:
				t.Fatalf("Wrong status code: got %v want %v. %s", code, tc.wantedCode, err)
			case code != 200:
				return
			}

			resB := &struct {
				ID        string `json="id"`
				Email     string `json="email"`
				Password  string `json="password"`
				FirstName string `json="firstName"`
				LastName  string `json="lastName"`
				Role      string `json="role"`
			}{}

			d := json.NewDecoder(w.Body)
			err = d.Decode(resB)
			if err != nil {
				t.Fatal(err)
			}

			switch {
			case resB.ID != u.ID:
				t.Errorf("Invalid parameters: %v or %v", resB.ID, u.ID)
			case resB.Email != u.Email:
				t.Errorf("Invalid parameters: %v or %v", resB.Email, u.Email)
			case resB.Password != "":
				t.Errorf("Password should be empty for security reasons")
			case resB.FirstName != u.FirstName:
				t.Errorf("Invalid parameters: %v or %v", resB.FirstName, u.FirstName)
			case resB.LastName != u.LastName:
				t.Errorf("Invalid parameters: %v or %v", resB.FirstName, u.FirstName)
			}
		})
	}
	s.Storage.DB.CleanupTables()
}

func TestReadAllUsers(t *testing.T) {
	ids := []string{"10cf445a-b18e-4f29-9139-57fab62aedbc", "20cf445a-b18e-4f29-9139-57fab62aedbc", "30cf445a-b18e-4f29-9139-57fab62aedbc"}
	ph := []byte("$2a$10$flStfMMZw4TsuJh3OdJhYeDBCibDlTNNm.yVMya4RgMcc7bF0/2nq")
	us := &[]rolf.User{
		{ID: ids[0], Email: "readAllUsers1@mail.com", Password: "secret", PasswordHash: ph, FirstName: "Rolf", LastName: "Baeckman", Role: "user"},
		{ID: ids[1], Email: "readAllUsers2@mail.com", Password: "secret", PasswordHash: ph, FirstName: "Rolf", LastName: "Baeckman", Role: "user"},
		{ID: ids[2], Email: "readAllUsers3@mail.com", Password: "secret", PasswordHash: ph, FirstName: "Rolf", LastName: "Baeckman", Role: "user"},
	}

	for _, u := range *us {
		err := s.Storage.DB.Create(&u)
		if err != nil {
			t.Fatal(err)
		}
	}

	r, err := http.NewRequest("GET", "localhost:3000/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := s.LogAndError(s.ReadAllUsers())

	w := httptest.NewRecorder()
	err = h(w, r)

	if code := w.Code; code != http.StatusOK {
		t.Fatalf("Wrong status code: got %v want %v. %s", code, http.StatusOK, err)
		return
	}

	resBs := &[]struct {
		ID        string `json="id"`
		Email     string `json="email"`
		Password  string `json="password"`
		FirstName string `json="firstName"`
		LastName  string `json="lastName"`
		Role      string `json="role"`
	}{}

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
	s.Storage.DB.CleanupTables()
}

func TestUpdateUser(t *testing.T) {
	// add two hex after id to make a whole uuid
	id := "10cf445a-b18e-4f29-9139-57fab62aed"
	tests := map[string]struct {
		ID         string `json="id,omitempty"`
		Email      string `json="email,omitempty"`
		Password   string `json="password,omitempty"`
		FirstName  string `json="firstName,omitempty"`
		LastName   string `json="lastName,omitempty"`
		wantedCode int    `json="-"`
	}{
		// url variable tests
		"id misformated": {ID: "$%)#(%'%')", Email: "updateUser1@mail.com", Password: "secret", FirstName: "UpdatedRolf", LastName: "Baeckman", wantedCode: 422},
		"id empty":       {ID: "", Email: "updateUser2@mail.com", Password: "secret", FirstName: "UpdatedRolf", LastName: "Baeckman", wantedCode: 422},

		// request body tests
		"simple":                 {ID: id + "01", Email: "updateUser4@mail.com", Password: "secret", FirstName: "UpdatedRolf", LastName: "Baeckman", wantedCode: 200},
		"email misformatted":     {ID: id + "02", Email: "updateUser5mail.com", Password: "secret", FirstName: "UpdatedRolf", LastName: "Baeckman", wantedCode: 422},
		"firstname misformatted": {ID: id + "03", Email: "updateUser6@mail.com", Password: "secret", FirstName: "UpdatedRolf#$", LastName: "Baeckman", wantedCode: 422},
		"lastname misformatted":  {ID: id + "04", Email: "updateUser7@mail.com", Password: "secret", FirstName: "UpdatedRolf", LastName: "Baeckman$%", wantedCode: 422},
		"password to short":      {ID: id + "05", Email: "updateUser8@mail.com", Password: "sec", FirstName: "UpdatedRolf", LastName: "Baeckman", wantedCode: 422},
		"email missing":          {ID: id + "06", Password: "updateUser9@mail.com", FirstName: "UpdatedRolf", LastName: "Baeckman", wantedCode: 422},
		"password missing":       {ID: id + "07", Email: "updateUser10@mail.com", FirstName: "UpdatedRolf", LastName: "Baeckman", wantedCode: 422},
		"firstname missing":      {ID: id + "08", Email: "updateUser11@mail.com", Password: "secret", LastName: "Baeckman", wantedCode: 422},
		"lastname missing":       {ID: id + "09", Email: "updateUser12@mail.com", Password: "secret", FirstName: "UpdatedRolf", wantedCode: 422},
	}

	for name, reqB := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(reqB)
			if err != nil {
				t.Fatal(err)
			}

			// TODO: change from static to env hostname
			r, err := http.NewRequest("PUT", "localhost:3000/v1/users", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}

			r = mux.SetURLVars(r, map[string]string{
				"id": reqB.ID,
			})

			h := s.LogAndError(s.UpdateUser())

			w := httptest.NewRecorder()
			err = h(w, r)

			switch code := w.Code; {
			case code != reqB.wantedCode:
				t.Fatalf("Wrong status code: got %v want %v", code, reqB.wantedCode)
			case code != 200:
				return
			}

			resB := &struct {
				ID        string `json="id"`
				Email     string `json="email"`
				Password  string `json="password"`
				FirstName string `json="firstName"`
				LastName  string `json="lastName"`
				Role      string `json="role"`
			}{}

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
				t.Errorf("Invalid parameters: %v or %v", resB.LastName, reqB.LastName)
			case resB.Role == "":
				t.Errorf("Invalid role: %v", resB.Role)
			}
		})
	}

	s.Storage.DB.CleanupTables()
}

func TestDeleteUser(t *testing.T) {
	id := "10ff445a-b18e-4f29-9139-57fab62aedbc"
	tests := map[string]struct {
		id         string
		wantedCode int `json="-"`
	}{
		"simple":         {id: id, wantedCode: 200},
		"id misformated": {id: "$%)#(%'%')", wantedCode: 422},
		"id empty":       {id: "", wantedCode: 422},
	}

	ph := []byte("$2a$10$flStfMMZw4TsuJh3OdJhYeDBCibDlTNNm.yVMya4RgMcc7bF0/2nq")
	u := &rolf.User{ID: id, Email: "deleteUser@mail.com", Password: "secret", PasswordHash: ph, FirstName: "Rolf", LastName: "Baeckman", Role: "user"}

	err := s.Storage.DB.Create(u)
	if err != nil {
		t.Fatal(err)
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := http.NewRequest("DELETE", "localhost:3000/v1/users", nil)
			if err != nil {
				t.Fatal(err)
			}

			r = mux.SetURLVars(r, map[string]string{
				"id": tc.id,
			})

			h := s.LogAndError(s.DeleteUser())

			w := httptest.NewRecorder()
			err = h(w, r)

			if code := w.Code; code != tc.wantedCode {
				t.Fatalf("Wrong status code: got %v want %v. %s", code, tc.wantedCode, err)
			}
		})
	}
	s.Storage.DB.CleanupTables()
}

func TestLogin(t *testing.T) {
	tests := map[string]struct {
		Email      string `json="email,omitempty"`
		Password   string `json="password,omitempty"`
		wantedCode int    `json="-"`
	}{
		"simple":             {Email: "loginUser@mail.com", Password: "secret", wantedCode: 200},
		"email misformatted": {Email: "loginUsermail.com", Password: "secret", wantedCode: 422},
		"email missing":      {Password: "secret", wantedCode: 422},
		"password missing":   {Email: "loginUser@mail.com", wantedCode: 422},
		//"sql injection":   {Email: "user8@mail.com", Password: "secret", wantedCode: 422},
	}

	id := "123f445a-b18e-4f29-9139-57fab62aedbc"
	ph := []byte("$2a$10$flStfMMZw4TsuJh3OdJhYeDBCibDlTNNm.yVMya4RgMcc7bF0/2nq")
	u := &rolf.User{ID: id, Email: "loginUser@mail.com", Password: "secret", PasswordHash: ph, FirstName: "Rolf", LastName: "Baeckman", Role: "user"}

	err := s.Storage.DB.Create(u)
	if err != nil {
		t.Fatal(err)
	}

	for name, reqB := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(reqB)
			if err != nil {
				t.Fatal(err)
			}

			// TODO: change from static to env hostname
			r, err := http.NewRequest("POST", "localhost:3000/v1/login", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}

			h := s.LogAndError(s.Login())

			w := httptest.NewRecorder()
			err = h(w, r)

			//c, err := r.Cookie(rolf.Env["COOKIE_NAME"])
			//if err != nil {
			//	return &rolf.Error{Code: rolf.EFORBIDDEN, Message: "You have to login first", Op: rolf.Trace(), Err: err}
			//	return err
			//}

			switch code := w.Code; {
			case code != reqB.wantedCode:
				t.Fatalf("Wrong status code: got %v want %v", code, reqB.wantedCode)
			case code != 200:
				return
			}

			resB := &struct {
				ID        string `json="id"`
				Email     string `json="email"`
				FirstName string `json="firstName"`
				LastName  string `json="lastName"`
				Role      string `json="role"`
			}{}

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
			case resB.Role == "":
				t.Errorf("Invalid role: %v", resB.Role)
			}
		})
	}

	s.Storage.DB.CleanupTables()
}

func BenchmarkCreateUser(b *testing.B) {
	reqB := &rolf.User{}
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

func generateUserData(u *rolf.User) {
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
func createRandomUsers(i int) (*[]rolf.User, error) {
	us := []rolf.User{}
	for y := 0; y < i; y++ {
		u := &rolf.User{}
		generateUserData(u)
		err := s.Storage.DB.Create(u)
		if err != nil {
			return &us, err
		}
		us = append(us, *u)
	}

	return &us, nil
}
