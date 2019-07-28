package rest

import (
	"fmt"
	"log"
	"os"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
	"github.com/RobinBaeckman/rolf/pkg/storage/memcache"
	"github.com/RobinBaeckman/rolf/pkg/storage/mysql"
	"github.com/RobinBaeckman/rolf/pkg/storage/redis"
	"github.com/gorilla/mux"
)

func NewServer() (*Server, error) {
	l := log.New(logWriter{}, "", 3)
	l.SetFlags(0)
	mdb := redis.NewMemDB()
	db, err := mysql.NewDB()
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	c := memcache.NewStorage(db, mdb)

	return &Server{
		Router:  r,
		Storage: c,
		Logger:  l,
	}, nil
}

func Run() error {
	if err := ParseEnv(); err != nil {
		return err
	}

	s, err := NewServer()
	defer s.Storage.DB.Close()
	if err != nil {
		return err
	}

	s.Routes()

	return nil
}

func ParseEnv() error {
	for key, _ := range rolf.Env {
		if v, ok := os.LookupEnv(key); ok {
			rolf.Env[key] = v
		} else {
			return fmt.Errorf("missing env variable: %s\n", key)
		}
	}

	return nil
}
