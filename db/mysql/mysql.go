package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/RobinBaeckman/ragnar"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_PASS")+"@tcp("+os.Getenv("MYSQL_HOST")+")/"+os.Getenv("MYSQL_DB"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type UserService struct {
	*sql.DB
}

func (s *UserService) Read(u *ragnar.User) error {
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

func (s *UserService) ReadByEmail(e string) (*ragnar.User, error) {
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

func (s *UserService) Create(u *ragnar.User) error {
	stmtIns, err := s.Prepare("INSERT INTO users(id, email, password, first_name, last_name, role) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(u.ID, u.Email, u.Password, u.FirstName, u.LastName, "user")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	return nil
}

func (s *UserService) ReadAll(us *[]ragnar.User) error {
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
