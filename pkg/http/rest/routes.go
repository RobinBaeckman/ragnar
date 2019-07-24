package rest

import (
	"log"
	"net/http"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
	"github.com/RobinBaeckman/ragnar/pkg/storage/memcache"
	"github.com/gorilla/mux"
)

// TODO: make sure I'm using the correct fields as exported.
type Server struct {
	Storage *memcache.Storage
	Router  *mux.Router
	Logger  *log.Logger
}

func (s *Server) Routes() {
	// Create
	s.Router.Handle("/v1/users",
		s.LogAndError(s.CreateUser())).Methods("POST")

	// Read
	s.Router.Handle("/v1/users/{id}",
		s.LogAndError(s.Auth(s.ReadUser()))).Methods("GET")

	// ReadAll
	s.Router.Handle("/v1/users",
		s.LogAndError(s.Auth(s.ReadAllUsers()))).Methods("GET")

	// Update
	s.Router.Handle("/v1/users/{id}",
		s.LogAndError(s.Auth(s.UpdateUser()))).Methods("PUT")

	// Delete
	s.Router.Handle("/v1/users/{id}",
		s.LogAndError(s.Auth(s.DeleteUser()))).Methods("DELETE")

	// Auth
	s.Router.Handle("/v1/login",
		s.LogAndError(s.Login())).Methods("POST")

	s.Router.Handle("/v1/logout",
		s.LogAndError(s.Auth(s.Logout()))).Methods("GET")

	http.Handle("/", s.Router)

	s.Logger.Printf("Running on: %s:%s", ragnar.Env["HOST"], ragnar.Env["PORT"])
	s.Logger.Fatal(http.ListenAndServe(ragnar.Env["HOST"]+":"+ragnar.Env["PORT"], nil))
}
