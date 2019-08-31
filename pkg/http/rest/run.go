package rest

import (
	"fmt"
	"log"
	"os"

	"github.com/RobinBaeckman/rolf/pkg/mail/smtp"
	"github.com/RobinBaeckman/rolf/pkg/rolf"
	"github.com/RobinBaeckman/rolf/pkg/storage/memcache"
	"github.com/RobinBaeckman/rolf/pkg/storage/mysql"
	"github.com/RobinBaeckman/rolf/pkg/storage/redis"
	"github.com/gorilla/mux"
)

func NewServer(l *log.Logger, mdb rolf.MemDB, db rolf.DB, m rolf.Mailer, r *mux.Router) (*Server, error) {
	c := memcache.NewStorage(db, mdb)

	return &Server{
		Router:  r,
		Storage: c,
		Logger:  l,
		Mailer:  m,
	}, nil
}

func Run() error {
	if err := ParseEnv(); err != nil {
		return err
	}

	l := log.New(logWriter{}, "", 3)
	l.SetFlags(0)
	mdb := redis.NewMemDB()
	db, err := mysql.NewDB()
	if err != nil {
		return err
	}
	defer db.Close()

	m, err := smtp.NewMailer()
	if err != nil {
		return err
	}

	r := mux.NewRouter()

	s, err := NewServer(l, mdb, db, m, r)
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
