package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
	"github.com/RobinBaeckman/ragnar/pkg/storage/memcache"
	"github.com/RobinBaeckman/ragnar/pkg/storage/mysql"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func Run() error {
	l := log.New(os.Stdout, "", 3)

	if err := parseEnv(); err != nil {
		return err
	}

	l.SetPrefix(ragnar.Env["LOG_PREFIX"])

	re := newRedis()
	db, err := mysql.NewDB()
	if err != nil {
		return err
	}
	defer db.Close()

	r := mux.NewRouter()

	mdb := &mysql.DB{db}
	c := memcache.NewUserStorage(mdb, re)

	s := &Server{
		router:      r,
		userStorage: c,
		logger:      l,
	}

	s.Routes()

	l.Printf("Running on: %s:%s", ragnar.Env["HOST"], ragnar.Env["PORT"])
	l.Fatal(http.ListenAndServe(ragnar.Env["HOST"]+":"+ragnar.Env["PORT"], nil))

	return nil
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
			return fmt.Errorf("missing env variable: %s\n", key)
		}
	}

	return nil
}
