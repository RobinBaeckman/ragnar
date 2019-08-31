package rest

import (
	"log"
	"net/http"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
	"github.com/RobinBaeckman/rolf/pkg/storage/memcache"
	"github.com/gorilla/mux"
)

// TODO: make sure I'm using the correct fields as exported.
type Server struct {
	Storage *memcache.Storage
	Router  *mux.Router
	Logger  *log.Logger
	Mailer  rolf.Mailer
}

func (s *Server) Routes() {
	// User
	s.Router.Handle("/v1/users/{id}",
		s.LogAndError(s.Auth(s.ReadUser()))).Methods("GET")

	s.Router.Handle("/v1/users/{id}",
		s.LogAndError(s.Auth(s.UpdateUser()))).Methods("PUT")

	s.Router.Handle("/v1/users/{id}",
		s.LogAndError(s.Auth(s.DeleteUser()))).Methods("DELETE")

	// Admin
	s.Router.Handle("/v1/admin/register",
		s.LogAndError(s.Auth(s.OnlyAdmin(s.RegisterAdmin())))).Methods("POST")
	s.Router.Handle("/v1/admin/users",
		s.LogAndError(s.Auth(s.OnlyAdmin(s.ReadAllUsers())))).Methods("GET")
	s.Router.Handle("/v1/admin/users/{id}",
		s.LogAndError(s.Auth(s.OnlyAdmin(s.ReadAnyUser())))).Methods("GET")

	// Login
	s.Router.Handle("/v1/register",
		s.LogAndError(s.RegisterUser())).Methods("POST")

	s.Router.Handle("/v1/users/{email}/forgot_password",
		s.LogAndError(s.ForgotPasswordMail())).Methods("GET")

	s.Router.Handle("/v1/users/{email}/forgot_password",
		s.LogAndError(s.PasswordReset())).Methods("POST")

	s.Router.Handle("/v1/login",
		s.LogAndError(s.Login())).Methods("POST")

	s.Router.Handle("/v1/refresh",
		s.LogAndError(s.Refresh())).Methods("POST")

	http.Handle("/", s.Router)

	s.Logger.Printf("Running on: %s:%s", rolf.Env["HOST"], rolf.Env["PORT"])
	s.Logger.Fatal(http.ListenAndServe(rolf.Env["HOST"]+":"+rolf.Env["PORT"], nil))
}
