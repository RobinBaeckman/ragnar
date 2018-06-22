package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RobinBaeckman/ragnar"
	"github.com/spf13/viper"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", viper.GetString("mysql.user")+":"+viper.GetString("mysql.password")+"@tcp("+viper.GetString("mysql.host")+")/"+viper.GetString("mysql.db"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type UserService struct {
	*sql.DB
}

func (s *UserService) Get(u *ragnar.User) error {
	err := s.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE id=?", u.ID).Scan(&u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Email is %s\n", u.Email)
	}

	return nil
}

func (s *UserService) GetByEmail(e string) (*ragnar.User, error) {
	u := &ragnar.User{}
	err := s.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE email=?", e).Scan(&u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Email is %s\n", u.Email)
	}

	return u, nil
}

func (s *UserService) Store(u *ragnar.User) error {
	stmtIns, err := s.Prepare("INSERT INTO users(id, email, password, first_name, last_name, role) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	stmtIns.Exec(u.ID, u.Email, u.Password, u.FirstName, u.LastName, "user")
	defer stmtIns.Close()

	return nil
}

func (s *UserService) GetAll(us *[]ragnar.User) error {
	rows, err := s.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for rows.Next() {
		u := ragnar.User{}
		err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		*us = append(*us, u)
	}

	return nil
}
