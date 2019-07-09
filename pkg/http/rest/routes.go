package rest

import (
	"log"
	"net/http"

	"github.com/RobinBaeckman/ragnar/pkg/storage/memcache"
	"github.com/gorilla/mux"
)

type Server struct {
	userStorage *memcache.UserStorage
	router      *mux.Router
	logger      *log.Logger
}

func (s *Server) Routes() {
	// Create
	s.router.Handle("/v1/users",
		s.checkError(s.log(s.CreateUser()))).Methods("POST")

	// Read
	s.router.Handle("/v1/users/{id}",
		s.checkError(s.log(s.auth(s.ReadUser())))).Methods("GET")

	// ReadAll
	s.router.Handle("/v1/users",
		s.checkError(s.log(s.auth(s.ReadAllUsers())))).Methods("GET")

	// Update
	s.router.Handle("/v1/users/{id}",
		s.checkError(s.log(s.auth(s.UpdateUser())))).Methods("PUT")

	// Delete
	s.router.Handle("/v1/users/{id}",
		s.checkError(s.log(s.auth(s.DeleteUser())))).Methods("DELETE")

	// Auth
	s.router.Handle("/v1/login",
		s.checkError(s.log(s.Login()))).Methods("POST")

	s.router.Handle("/v1/logout",
		s.checkError(s.log(s.auth(s.Logout())))).Methods("GET")

	http.Handle("/", s.router)
}
