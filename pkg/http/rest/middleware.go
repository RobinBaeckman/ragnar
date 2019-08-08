package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte(rolf.Env["JWT_KEY"])

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
			if err == http.ErrNoCookie {
				return &rolf.Error{Code: http.StatusForbidden, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
			}

			return &rolf.Error{Code: http.StatusForbidden, Msg: "You have to login first", Op: rolf.Trace() + err.Error()}
		}

		tknStr := c.Value

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

		ctx := context.WithValue(r.Context(), "myID", cl.Subject)
		ctx = context.WithValue(ctx, "role", cl.Role)
		err = next(w, r.WithContext(ctx))
		if err != nil {
			switch v := err.(type) {
			case *rolf.Error:
				// TODO: Make user part of the error instead
				v.Op = fmt.Sprintf("%s User: %s", v.Op, cl.Subject)
			}
			return err
		}

		return nil
	}
}

func (s *Server) IsAdmin(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		if role := r.Context().Value("role"); role != "admin" {
			return &rolf.Error{Code: http.StatusUnauthorized, Msg: "you are not authorized because you are not admin", Op: rolf.Trace()}
		}
		return nil
	}
}

func (s *Server) LogAndError(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		logOutput := fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		s.Logger.Printf(logOutput)
		if err := next(w, r); err != nil {
			v := err.(*rolf.Error)
			http.Error(w, v.Msg, v.Code)
		}

		return err
	}
}
