package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", ragnar.Env["MYSQL_USER"]+":"+ragnar.Env["MYSQL_PASS"]+"@tcp("+ragnar.Env["MYSQL_HOST"]+")/"+ragnar.Env["MYSQL_DB"])
	if err != nil {
		return nil, fmt.Errorf("Can't connect to db\n")
	}

	return db, nil
}

type DB struct {
	*sql.DB
}

func (s *DB) Create(u *ragnar.User) error {
	stmtIns, err := s.Prepare("INSERT INTO users(id, email, password, first_name, last_name, role) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return &ragnar.Error{Code: ragnar.ECONFLICT, Message: "Username already exists", Op: ragnar.Trace(), Err: err}
	}

	_, err = stmtIns.Exec(u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role)
	if err != nil {
		return &ragnar.Error{Code: ragnar.ECONFLICT, Message: "Username already exists", Op: ragnar.Trace(), Err: err}
	}
	defer stmtIns.Close()

	return nil
}

func (s *DB) Read(u *ragnar.User) error {
	err := s.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE id=?", u.ID).Scan(&u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: "No user with that ID.", Op: ragnar.Trace(), Err: err}
	case err != nil:
		log.Fatal(err)
	}

	return nil
}

func (s *DB) ReadByEmail(u *ragnar.User) error {
	err := s.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE email=?", u.Email).Scan(&u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: "No user with that ID.", Op: ragnar.Trace(), Err: err}
	case err != nil:
		log.Fatal(err)
		return err
	default:
		fmt.Printf("Email is %s\n", u.Email)
	}

	return nil
}

func (s *DB) ReadAll(us *[]ragnar.User) error {
	rows, err := s.Query("SELECT * FROM users")
	if err != nil {
		return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: "No user with that ID.", Op: ragnar.Trace(), Err: err}
	}

	for rows.Next() {
		u := ragnar.User{}
		err = rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt)
		if err != nil {
			return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: "No user with that ID.", Op: ragnar.Trace(), Err: err}
		}
		*us = append(*us, u)
	}

	return nil
}

// TODO: fix role system
// TODO: duplication of password hash generation
func (s *DB) Update(u *ragnar.User) error {
	stmtIns, err := s.Prepare("UPDATE users set email=?, password=?, first_name=?, last_name=?, role=? where id=?")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role, u.ID)
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	return nil
}

func (s *DB) Delete(u *ragnar.User) error {
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
