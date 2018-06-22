package mongo

import (
	"github.com/RobinBaeckman/atla/errors"
	"github.com/RobinBaeckman/ragnar"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func NewDB() *mgo.Database {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("atla")

	return c
}

type UserService struct {
	Conn *mgo.Database
}

func (s *UserService) Get(a *ragnar.User) error {

	err := s.Conn.C("users").Find(bson.M{"id": a.ID}).One(a)
	if err != nil {
		return &errors.ErrHTTP{err, "There is no user user with that id.", 404}
	}

	return nil
}

func (s *UserService) Persist(a *ragnar.User) error {
	err := s.Conn.C("users").Insert(a)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetAll(as *[]ragnar.User) error {
	err := s.Conn.C("users").Find(nil).All(as)
	if err != nil {
		return err
	}

	return nil
}
