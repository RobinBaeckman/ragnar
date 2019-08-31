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

func (s *Server) RegisterUser() func(http.ResponseWriter, *http.Request) error {
	// TODO: thing := prepareThing()
	type request struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	type response struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) (err error) {
		// use thing
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		req := &request{}
		err = decoder.Decode(req)
		if err != nil {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Missing parameters", Op: rolf.Trace() + err.Error()}
		}

		switch {
		case !valid.Email(req.Email):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid email: " + req.Email, Op: rolf.Trace()}
		case !valid.Password(req.Password):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid password: " + req.Password, Op: rolf.Trace()}
		case !valid.FirstName(req.FirstName):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid firstName: " + req.FirstName, Op: rolf.Trace()}
		case !valid.LastName(req.LastName):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid lastName: " + req.LastName, Op: rolf.Trace()}
		}

		u := &rolf.User{
			ID:        uuid.New().String(),
			Email:     req.Email,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			// TODO: build better Role implementation
			Role: "user",
		}

		// TODO: make the hash in the database instead
		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		err = s.Storage.Create(u)
		if err != nil {
			return err
		}

		res := response{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		}

		b, err := json.Marshal(res)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		url := fmt.Sprintf("%s://%s:%s/v1/users/%s", rolf.Env["PROTO"], rolf.Env["HOST"], rolf.Env["PORT"], u.ID)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusCreated)
		w.Write(b)

		return nil
	}
}

func (s *Server) RegisterAdmin() func(http.ResponseWriter, *http.Request) error {
	// TODO: thing := prepareThing()
	type request struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	type response struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) (err error) {
		// use thing
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		req := &request{}
		err = decoder.Decode(req)
		if err != nil {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Missing parameters", Op: rolf.Trace() + err.Error()}
		}

		switch {
		case !valid.Email(req.Email):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid email: " + req.Email, Op: rolf.Trace()}
		case !valid.Password(req.Password):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid password: " + req.Password, Op: rolf.Trace()}
		case !valid.FirstName(req.FirstName):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid firstName: " + req.FirstName, Op: rolf.Trace()}
		case !valid.LastName(req.LastName):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid lastName: " + req.LastName, Op: rolf.Trace()}
		}

		u := &rolf.User{
			ID:        uuid.New().String(),
			Email:     req.Email,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			// TODO: build better Role implementation
			Role: "admin",
		}

		// TODO: make the hash in the database instead
		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		err = s.Storage.Create(u)
		if err != nil {
			return err
		}

		res := response{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		}

		b, err := json.Marshal(res)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
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
	type request struct {
		ID string `json:"id"`
	}
	type response struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		req := request{ID: mux.Vars(r)["id"]}

		if myID := r.Context().Value("myID"); myID != req.ID {
			return &rolf.Error{Code: http.StatusForbidden, Msg: "You don't have permission to view this resource", Op: rolf.Trace()}
		}

		if !valid.UUID(req.ID) {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid UUID.", Op: rolf.Trace()}
		}

		u := &rolf.User{ID: req.ID}
		err = s.Storage.Read(u)
		if err != nil {
			return err
		}

		res := response{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		}

		b, err := json.Marshal(res)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) ReadAnyUser() func(http.ResponseWriter, *http.Request) error {
	type request struct {
		ID string `json:"id"`
	}
	type response struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		req := request{ID: mux.Vars(r)["id"]}

		if !valid.UUID(req.ID) {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid UUID.", Op: rolf.Trace()}
		}

		u := &rolf.User{ID: req.ID}
		err = s.Storage.ReadAny(u)
		if err != nil {
			return err
		}

		res := response{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		}

		b, err := json.Marshal(res)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

// TODO: make sure the password is returned too
func (s *Server) ReadAllUsers() func(http.ResponseWriter, *http.Request) error {
	type response struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		us := &[]rolf.User{}
		// TODO: add ReadAll to memcache
		err = s.Storage.DB.ReadAll(us)
		if err != nil {
			return err
		}

		ress := []response{}
		for _, u := range *us {
			res := response{
				ID:        u.ID,
				Email:     u.Email,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Role:      u.Role,
			}
			ress = append(ress, res)
		}

		// Todo: create a map called mapUserToResp
		b, err := json.Marshal(ress)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) UpdateUser() func(http.ResponseWriter, *http.Request) error {
	type request struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	type response struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		req := &request{ID: mux.Vars(r)["id"]}
		err = decoder.Decode(req)
		if err != nil {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Missing parameters", Op: rolf.Trace() + err.Error()}
		}

		// TODO: integrate with validate function
		if !valid.UUID(req.ID) {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid UUID: " + req.ID, Op: rolf.Trace()}
		}

		switch {
		case !valid.Email(req.Email):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid email: " + req.Email, Op: rolf.Trace()}
		case !valid.Password(req.Password):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid password: " + req.Password, Op: rolf.Trace()}
		case !valid.FirstName(req.FirstName):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid firstName: " + req.FirstName, Op: rolf.Trace()}
		case !valid.LastName(req.LastName):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid lastName: " + req.LastName, Op: rolf.Trace()}
		}

		u := &rolf.User{
			ID:        req.ID,
			Email:     req.Email,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			// TODO: build better Role implementation
			Role: "user",
		}

		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		err = s.Storage.Update(u)
		if err != nil {
			return err
		}

		res := response{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		}

		b, err := json.Marshal(res)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) DeleteUser() func(http.ResponseWriter, *http.Request) error {
	type request struct {
		ID string `json:"id"`
	}
	type response struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) (err error) {

		req := request{ID: mux.Vars(r)["id"]}

		if !valid.UUID(req.ID) {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid UUID: " + req.ID, Op: rolf.Trace()}
		}

		u := &rolf.User{ID: req.ID}
		err = s.Storage.Delete(u)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}

func (s *Server) Login() func(http.ResponseWriter, *http.Request) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		req := &request{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		err = decoder.Decode(req)
		if err != nil {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Missing parameters", Op: rolf.Trace() + err.Error()}
		}

		switch {
		case !valid.Email(req.Email):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid email: " + req.Email, Op: rolf.Trace()}
		case !valid.Password(req.Password):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid password: " + req.Password, Op: rolf.Trace()}
		}
		u := &rolf.User{Email: req.Email, Password: req.Password}

		if err = s.Storage.DB.ReadByEmail(u); err != nil {
			return err
		}

		if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(u.Password)); err != nil {
			return &rolf.Error{Code: http.StatusUnauthorized, Msg: "Wrong username or password.", Op: rolf.Trace() + err.Error()}
		}

		et := time.Now().Add(5 * time.Minute)
		claims := &claims{
			Role: u.Role,
			StandardClaims: jwt.StandardClaims{
				Subject:   u.ID,
				ExpiresAt: et.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		c := http.Cookie{
			Name:     rolf.Env["COOKIE_NAME"],
			Value:    tokenString,
			HttpOnly: true,
			Expires:  et,
		}
		http.SetCookie(w, &c)

		res := response{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		}

		b, err := json.Marshal(res)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

		return nil
	}
}

func (s *Server) Refresh() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		c, err := r.Cookie(rolf.Env["COOKIE_NAME"])
		if err != nil {
			if err == http.ErrNoCookie {
				return &rolf.Error{Code: http.StatusForbidden, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
			}

			return &rolf.Error{Code: http.StatusForbidden, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
		}

		tknStr := c.Value

		// Initialize a new instance of `Claims`
		claims := &claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if !tkn.Valid {
			return &rolf.Error{Code: http.StatusUnauthorized, Msg: "you are not authorized", Op: rolf.Trace()}
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return &rolf.Error{Code: http.StatusUnauthorized, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
			}
			return &rolf.Error{Code: http.StatusBadRequest, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return &rolf.Error{Code: http.StatusUnauthorized, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
			}
			return &rolf.Error{Code: http.StatusBadRequest, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
		}

		// We ensure that a new token is not issued until enough time has elapsed
		// In this case, a new token will only be issued if the old token is within
		// 30 seconds of expiry. Otherwise, return a bad request status
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
			return &rolf.Error{Code: http.StatusBadRequest, Msg: "something", Op: rolf.Trace()}
		}

		// Now, create a new token for the current use, with a renewed expiration time
		et := time.Now().Add(5 * time.Minute)
		claims.ExpiresAt = et.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		// Set the new token as the users `token` cookie
		http.SetCookie(w, &http.Cookie{
			Name:    rolf.Env["COOKIE_NAME"],
			Value:   tokenString,
			Expires: et,
		})

		return nil
	}
}

func (s *Server) ForgotPasswordMail() func(http.ResponseWriter, *http.Request) error {
	type request struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		req := request{Email: mux.Vars(r)["email"]}

		if err != nil {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Missing parameters", Op: rolf.Trace() + err.Error()}
		}

		if !valid.Email(req.Email) {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid email: " + req.Email, Op: rolf.Trace()}
		}

		u := &rolf.User{Email: req.Email}

		if err = s.Storage.DB.ReadByEmail(u); err != nil {
			return err
		}

		et := time.Now().Add(60 * time.Minute)
		claims := &claims{
			StandardClaims: jwt.StandardClaims{
				Subject:   u.ID,
				ExpiresAt: et.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		u.Token = tokenString
		if err = s.Storage.DB.StoreToken(u); err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		s.Mailer.Send(u.Email, rolf.Env["EMAIL"], tokenString)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		return nil
	}
}

func (s *Server) PasswordReset() func(http.ResponseWriter, *http.Request) error {
	type request struct {
		Password string `json:"password"`
		Token    string `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		email := mux.Vars(r)["email"]
		req := &request{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		err = decoder.Decode(req)
		if err != nil {
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Missing parameters", Op: rolf.Trace() + err.Error()}
		}

		switch {
		case !valid.Email(email):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid email: " + email, Op: rolf.Trace()}
		case !valid.Password(req.Password):
			return &rolf.Error{Code: http.StatusUnprocessableEntity, Msg: "Invalid password: " + req.Password, Op: rolf.Trace()}
		}

		tknStr := req.Token

		// Initialize a new instance of `Claims`
		cl := &claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, cl, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if !tkn.Valid {
			return &rolf.Error{Code: http.StatusUnauthorized, Msg: "you are not authorized", Op: rolf.Trace() + err.Error()}
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return &rolf.Error{Code: http.StatusUnauthorized, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
			}
			return &rolf.Error{Code: http.StatusBadRequest, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
		}

		u := &rolf.User{Email: email}
		// TODO: make the hash in the database instead
		s.Storage.DB.ReadByEmail(u)

		if u.Token != tknStr {
			return &rolf.Error{Code: http.StatusBadRequest, Msg: "You have to login first", Op: rolf.Trace()}
		}

		u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
		}

		s.Storage.DB.UpdatePassword(u)

		return nil
	}
}
