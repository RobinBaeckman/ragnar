// TODO: make sure all returned error values are correct and fix error stuff
package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/RobinBaeckman/rolf/pkg/rolf"
	_ "github.com/go-sql-driver/mysql"
)

// TODO: Implement new error handling
func NewDB() (DB, error) {
	db, err := sql.Open("mysql", rolf.Env["MYSQL_USER"]+":"+rolf.Env["MYSQL_PASS"]+"@tcp("+rolf.Env["MYSQL_HOST"]+")/"+rolf.Env["MYSQL_DB"])
	if err != nil {
		return DB{}, fmt.Errorf("Can't connect to db\n %s", err)
	}

	return DB{db}, nil
}

type DB struct {
	*sql.DB
}

func (db DB) Create(u *rolf.User) error {
	stmtIns, err := db.Prepare("INSERT INTO users(id, email, password, first_name, last_name, role) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return &rolf.Error{Code: http.StatusConflict, Op: rolf.Trace() + err.Error()}
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role)
	if err != nil {
		return &rolf.Error{Code: http.StatusConflict, Op: rolf.Trace() + err.Error()}
	}

	return nil
}

func (db DB) Read(u *rolf.User) error {
	err := db.QueryRow("SELECT email, first_name, last_name, role FROM users WHERE id=?", u.ID).Scan(&u.Email, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &rolf.Error{Code: http.StatusNotFound, Msg: "No user with that ID: " + u.ID, Op: rolf.Trace() + err.Error()}
	case err != nil:
		log.Fatal(err)
	}

	return nil
}

func (db DB) ReadAny(u *rolf.User) error {
	err := db.QueryRow("SELECT email, first_name, last_name, role FROM users WHERE id=?", u.ID).Scan(&u.Email, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &rolf.Error{Code: http.StatusNotFound, Msg: "No user with that ID: " + u.ID, Op: rolf.Trace() + err.Error()}
	case err != nil:
		log.Fatal(err)
	}

	return nil
}

func (db DB) ReadByEmail(u *rolf.User) error {
	err := db.QueryRow("SELECT id, email, password, first_name, last_name, role FROM users WHERE email=?", u.Email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &rolf.Error{Code: http.StatusNotFound, Msg: "No user with that ID: " + u.ID, Op: rolf.Trace() + err.Error()}
	case err != nil:
		log.Fatal(err)
		return err
	default:
	}

	return nil
}

func (db DB) ReadAll(us *[]rolf.User) error {
	rows, err := db.Query("SELECT id, email, first_name, last_name, role FROM users")
	if err != nil {
		return &rolf.Error{Code: http.StatusNotFound, Op: rolf.Trace() + err.Error()}
	}

	for rows.Next() {
		u := rolf.User{}
		err = rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role)
		if err != nil {
			return &rolf.Error{Code: http.StatusNotFound, Msg: "No user with that ID: " + u.ID, Op: rolf.Trace() + err.Error()}
		}
		*us = append(*us, u)
	}

	return nil
}

func (db DB) Update(u *rolf.User) error {
	stmtIns, err := db.Prepare("UPDATE users set email=?, password=?, first_name=?, last_name=?, role=? where id=?")
	defer stmtIns.Close()
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// TODO: maybe implement soft delete
func (db DB) Delete(u *rolf.User) error {
	stmtIns, err := db.Prepare("DELETE from users where id=?")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) Close() {
	defer db.DB.Close()
}

func (db DB) CleanupTables() error {
	_, err := db.Query("TRUNCATE TABLE users")
	if err != nil {
		return &rolf.Error{Code: http.StatusInternalServerError, Op: rolf.Trace() + err.Error()}
	}

	return nil
}
