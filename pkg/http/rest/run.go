package rest

import (
	"fmt"
	"log"
	"os"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
	"github.com/RobinBaeckman/ragnar/pkg/storage/memcache"
	"github.com/RobinBaeckman/ragnar/pkg/storage/mysql"
	"github.com/RobinBaeckman/ragnar/pkg/storage/redis"
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
	// TODO: Fix so db and other things that should close closes at the end of this function
	//defer s.Storage.DB.Close()

	s.Routes()

	return nil
}

func ParseEnv() error {
	for key, _ := range ragnar.Env {
		if v, ok := os.LookupEnv(key); ok {
			ragnar.Env[key] = v
		} else {
			return fmt.Errorf("missing env variable: %s\n", key)
		}
	}

	return nil
}
