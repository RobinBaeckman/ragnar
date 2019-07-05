package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/RobinBaeckman/ragnar"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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

func (s *UserService) Create(u *ragnar.User) error {
	buid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	uid := buid.String()

	pHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.ID = uid
	u.Role = "user"

	stmtIns, err := s.Prepare("INSERT INTO users(id, email, password, first_name, last_name, role) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(u.ID, u.Email, pHash, u.FirstName, u.LastName, u.Role)
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	return nil
}

func (s *UserService) Read(u *ragnar.User) error {
	err := s.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE id=?", u.ID).Scan(&u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
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

func (s *UserService) ReadByEmail(u *ragnar.User) error {
	err := s.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE email=?", u.Email).Scan(&u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
		return err
	case err != nil:
		log.Fatal(err)
		return err
	default:
		fmt.Printf("Email is %s\n", u.Email)
	}

	return nil
}

func (s *UserService) ReadAll(us *[]ragnar.User) error {
	rows, err := s.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for rows.Next() {
		u := ragnar.User{}
		err = rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		*us = append(*us, u)
	}

	return nil
}

// TODO: fix role system
// TODO: duplication of password hash generation
func (s *UserService) Update(u *ragnar.User) error {
	pHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Role = "user"

	stmtIns, err := s.Prepare("UPDATE users set email=?, password=?, first_name=?, last_name=?, role=? where id=?")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(u.Email, pHash, u.FirstName, u.LastName, "user", u.ID)
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	return nil
}

func (s *UserService) Delete(u *ragnar.User) error {
	stmtIns, err := s.Prepare("DELETE from users where id=?")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(u.ID)
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	return nil
}
