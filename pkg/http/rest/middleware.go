package rest

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
)

type HandlerFuncWithError func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFuncWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	logOutput := fmt.Sprintf("[%s  %s  ]\t%s", time.Now().UTC().Format("2006-01-02T15:04:05"), rolf.Env["LOG_PREFIX"], string(bytes))
	return fmt.Printf(logOutput)
}

func (s *Server) Auth(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		c, err := r.Cookie(rolf.Env["COOKIE_NAME"])
		if err != nil {
			return &rolf.Error{Code: rolf.EFORBIDDEN, Message: "You have to login first", Op: rolf.Trace(), Err: err}
			return err
		}

		email, err := s.Storage.MemDB.Get(c.Value)
		if err != nil {
			return &rolf.Error{Code: rolf.EFORBIDDEN, Message: "You have to login first", Op: rolf.Trace(), Err: err}
			return err
		}

		err = next(w, r)
		if err != nil {
			switch v := err.(type) {
			case *rolf.Error:
				// TODO: Make user part of the error instead
				v.Op = fmt.Sprintf("%s User: %s", v.Op, email)
			}
			return err
		}

		return nil
	}
}

func (s *Server) LogAndError(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		logOutput := fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		s.Logger.Printf(logOutput)
		if err := next(w, r); err != nil {
			switch v := err.(type) {
			case *rolf.Error:
				if v.Code == rolf.EUNAUTHORIZED {
					http.Error(w, v.Message, 401)
				} else if v.Code == rolf.ENOTFOUND {
					http.Error(w, v.Message, 404)
				} else if v.Code == rolf.EFORBIDDEN {
					http.Error(w, v.Message, 403)
				} else if v.Code == rolf.ECONFLICT {
					http.Error(w, v.Message, 409)
				} else if v.Code == rolf.EINVALID {
					http.Error(w, v.Message, 422)
				} else if v.Code == rolf.EINTERNAL {
					http.Error(w, rolf.EINTERNAL_MSG, 500)
				} else {
					http.Error(w, v.Error(), 500)
					logOutput = fmt.Sprintf("Status: %v, Error: %v", 500, v.Error())
					s.Logger.Printf(logOutput)
					debug.PrintStack()
					http.Error(w, v.Message, 500)
				}
				s.Logger.Println(err)
			}
			return err
		}

		return err
	}
}
