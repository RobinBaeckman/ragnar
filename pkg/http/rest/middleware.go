package rest

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
)

type HandlerFuncWithError func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFuncWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

type logWriter struct {
}

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func (writer logWriter) Write(bytes []byte) (int, error) {
	logOutput := fmt.Sprintf("[%s  %s  ]\t%s", time.Now().UTC().Format("2006-01-02T15:04:05"), ragnar.Env["LOG_PREFIX"], string(bytes))
	return fmt.Printf(logOutput)
}

func (s *Server) Auth(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		c, err := r.Cookie(ragnar.Env["COOKIE_NAME"])
		if err != nil {
			return &ragnar.Error{Code: ragnar.EFORBIDDEN, Message: "You have to login first", Op: ragnar.Trace(), Err: err}
			return err
		}

		email, err := s.Storage.MemDB.Get(c.Value)
		if err != nil {
			return &ragnar.Error{Code: ragnar.EFORBIDDEN, Message: "You have to login first", Op: ragnar.Trace(), Err: err}
			return err
		}

		err = next(w, r)
		if err != nil {
			switch v := err.(type) {
			case *ragnar.Error:
				// TODO: Make user part of the error instead
				v.Op = fmt.Sprintf("%s User: %s", v.Op, email)
			}
			return err
		}

		return nil
	}
}

func (s *Server) Log(next HandlerFuncWithError) HandlerFuncWithError {
	textColor := InfoColor
	return func(w http.ResponseWriter, r *http.Request) error {
		logOutput := fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		s.Logger.Printf(textColor, logOutput)
		err := next(w, r)
		if err != nil {
			return err
		}

		if textColor == InfoColor {
			textColor = WarningColor
		} else {
			textColor = InfoColor
		}

		return nil
	}
}

func (s *Server) CheckError(next HandlerFuncWithError) HandlerFuncWithError {
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
					s.Logger.Print(v)
				} else {
					s.Logger.Print(v)
				}
			default:
				http.Error(w, v.Error(), 500)
				s.Logger.Printf("Status: %v, Error: %v", 500, v.Error())
				debug.PrintStack()

			}
		}
		return err
	}
}
