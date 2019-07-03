package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RobinBaeckman/ragnar"
	"github.com/RobinBaeckman/ragnar/controller"
	"github.com/RobinBaeckman/ragnar/db/mysql"
	"github.com/RobinBaeckman/ragnar/errors"
	"github.com/RobinBaeckman/ragnar/middleware"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, os.Getenv("LOG_PREFIX"), 3)
	re := newRedis()
	db := mysql.NewDB()
	defer db.Close()

	r := mux.NewRouter()

	s := &mysql.UserService{db}
	c := ragnar.NewUserCache(s, re)

	// User user
	r.Handle("/v1/users/signup", Adapt(
		errors.Check(controller.CreateUser(c)),
		middleware.Log(l),
	)).Methods("POST")

	r.Handle("/v1/users/{id}", Adapt(
		errors.Check(controller.ReadUser(c)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("GET")

	r.Handle("/v1/users", Adapt(
		errors.Check(controller.ReadAllUsers(s)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("GET")

	r.Handle("/v1/login", Adapt(
		errors.Check(controller.Login(c)),
		middleware.Log(l),
	)).Methods("POST")

	r.Handle("/v1/logout", Adapt(
		errors.Check(controller.Logout(re)),
		middleware.Log(l),
		middleware.Auth(re),
	)).Methods("GET")

	http.Handle("/", r)

	l.Printf("Running on: %s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	l.Fatal(http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), nil))
}

func Adapt(h http.Handler, adapters ...middleware.Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func newRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
