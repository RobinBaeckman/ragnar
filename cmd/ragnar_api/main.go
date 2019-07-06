package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RobinBaeckman/ragnar"
	"github.com/RobinBaeckman/ragnar/db/mysql"
	"github.com/RobinBaeckman/ragnar/errors"
	"github.com/RobinBaeckman/ragnar/handler"
	"github.com/RobinBaeckman/ragnar/middleware"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "", 3)

	if err := parseEnv(); err != nil {
		l.Println(err)
		return
	}

	l.SetPrefix(ragnar.Env["LOG_PREFIX"])

	re := newRedis()
	db, err := mysql.NewDB()
	if err != nil {
		l.Println(err)
		return
	}
	defer db.Close()

	r := mux.NewRouter()

	s := &mysql.UserService{db}
	c := ragnar.NewUserCache(s, re)

	// Create
	r.Handle("/v1/users", Adapt(
		errors.Check(handler.CreateUser(c)),
		middleware.Log(l),
	)).Methods("POST")

	// Read
	r.Handle("/v1/users/{id}", Adapt(
		errors.Check(handler.ReadUser(c)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("GET")

	// ReadAll
	r.Handle("/v1/users", Adapt(
		errors.Check(handler.ReadAllUsers(s)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("GET")

	// Update
	r.Handle("/v1/users/{id}", Adapt(
		errors.Check(handler.UpdateUser(c)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("PUT")

	// Delete
	r.Handle("/v1/users/{id}", Adapt(
		errors.Check(handler.DeleteUser(c)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("DELETE")

	// Auth
	r.Handle("/v1/login", Adapt(
		errors.Check(handler.Login(c)),
		middleware.Log(l),
	)).Methods("POST")

	r.Handle("/v1/logout", Adapt(
		errors.Check(handler.Logout(re)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("GET")

	http.Handle("/", r)

	l.Printf("Running on: %s:%s", ragnar.Env["HOST"], ragnar.Env["PORT"])
	l.Fatal(http.ListenAndServe(ragnar.Env["HOST"]+":"+ragnar.Env["PORT"], nil))

	return
}

// TODO: make it so the handlers and middleware are read in the correct order in Adapt
func Adapt(h http.Handler, adapters ...middleware.Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func newRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     ragnar.Env["REDIS_HOST"] + ":" + ragnar.Env["REDIS_PORT"],
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func parseEnv() error {
	for key, _ := range ragnar.Env {
		if v, ok := os.LookupEnv(key); ok {
			ragnar.Env[key] = v
		} else {
			return fmt.Errorf("missing env variable: %s", key)
		}
	}

	return nil
}
