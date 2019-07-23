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
		s.CheckError(s.Log(s.CreateUser()))).Methods("POST")

	// Read
	s.Router.Handle("/v1/users/{id}",
		s.CheckError(s.Log(s.Auth(s.ReadUser())))).Methods("GET")

	// ReadAll
	s.Router.Handle("/v1/users",
		s.CheckError(s.Log(s.Auth(s.ReadAllUsers())))).Methods("GET")

	// Update
	s.Router.Handle("/v1/users/{id}",
		s.CheckError(s.Log(s.Auth(s.UpdateUser())))).Methods("PUT")

	// Delete
	s.Router.Handle("/v1/users/{id}",
		s.CheckError(s.Log(s.Auth(s.DeleteUser())))).Methods("DELETE")

	// Auth
	s.Router.Handle("/v1/login",
		s.CheckError(s.Log(s.Login()))).Methods("POST")

	s.Router.Handle("/v1/logout",
		s.CheckError(s.Log(s.Auth(s.Logout())))).Methods("GET")

	http.Handle("/", s.Router)

	s.Logger.Printf("Running on: %s:%s", ragnar.Env["HOST"], ragnar.Env["PORT"])
	s.Logger.Fatal(http.ListenAndServe(ragnar.Env["HOST"]+":"+ragnar.Env["PORT"], nil))
}
