package rest

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
)

type HandlerFuncWithError func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFuncWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func (s *Server) auth(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		c, err := r.Cookie(os.Getenv("COOKIE_NAME"))
		if err != nil {
			return &ragnar.Error{Code: ragnar.EFORBIDDEN, Message: "You have to login first", Op: ragnar.Trace(), Err: err}
			return err
		}

		email, err := s.userStorage.Redis.Get(c.Value).Result()
		if err != nil {
			return &ragnar.Error{Code: ragnar.EFORBIDDEN, Message: "You have to login first", Op: ragnar.Trace(), Err: err}
			return err
		}

		err = next(w, r)
		if err != nil {
			switch v := err.(type) {
			case *ragnar.Error:
				v.Op = fmt.Sprintf("%s #####User: %s", v.Op, email)
			}
			return err
		}

		return nil
	}
}

func (s *Server) log(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		s.logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		err := next(w, r)
		if err != nil {
			return err
		}

		return nil
	}
}

func (s *Server) checkError(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		if err := next(w, r); err != nil {
			switch v := err.(type) {
			case *ragnar.Error:
				if v.Code == ragnar.EUNAUTHORIZED {
					http.Error(w, v.Message, 401)
				} else if v.Code == ragnar.ENOTFOUND {
					http.Error(w, v.Message, 404)
				} else if v.Code == ragnar.EFORBIDDEN {
					http.Error(w, v.Message, 403)
				} else if v.Code == ragnar.ECONFLICT {
					http.Error(w, v.Message, 409)
				} else if v.Code == ragnar.EINVALID {
					http.Error(w, v.Message, 422)
				} else if v.Code == ragnar.EINTERNAL {
					http.Error(w, ragnar.EINTERNAL_MSG, 500)
				} else {
					http.Error(w, v.Message, 500)
				}

				if v.Err != nil {
					s.logger.Printf("#####Code: %v, #####Message: %v, #####Op: %v, #####Error: %v", v.Code, v.Message, v.Op, v.Err)
				} else {
					s.logger.Printf("#####Status: %v, #####Message: %v, #####Op: %v", v.Code, v.Message, v.Op)
				}
			default:
				http.Error(w, v.Error(), 500)
				s.logger.Printf("Status: %v, Error: %v", 500, v.Error())
				debug.PrintStack()

			}
		}
		return err
	}
}
