// TODO: make sure all returned error values are correct and fix error stuff
package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RobinBaeckman/ragnar/pkg/ragnar"
	_ "github.com/go-sql-driver/mysql"
)

// TODO: Implement new error handling
func NewDB() (DB, error) {
	db, err := sql.Open("mysql", ragnar.Env["MYSQL_USER"]+":"+ragnar.Env["MYSQL_PASS"]+"@tcp("+ragnar.Env["MYSQL_HOST"]+")/"+ragnar.Env["MYSQL_DB"])
	if err != nil {
		return DB{}, fmt.Errorf("Can't connect to db\n %s", err)
	}

	return DB{db}, nil
}

type DB struct {
	*sql.DB
}

func (db DB) Create(u *ragnar.User) error {
	stmtIns, err := db.Prepare("INSERT INTO users(id, email, password, first_name, last_name, role) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return &ragnar.Error{Code: ragnar.ECONFLICT, Op: ragnar.Trace(), Err: err}
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role)
	if err != nil {
		return &ragnar.Error{Code: ragnar.ECONFLICT, Op: ragnar.Trace(), Err: err}
	}

	return nil
}

func (db DB) Read(u *ragnar.User) error {
	err := db.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE id=?", u.ID).Scan(&u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: fmt.Sprintf("No user with that ID: %s", u.ID), Op: ragnar.Trace(), Err: err}
	case err != nil:
		log.Fatal(err)
	}

	return nil
}

func (db DB) ReadByEmail(u *ragnar.User) error {
	err := db.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE email=?", u.Email).Scan(&u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: fmt.Sprintf("No user with that ID: %s", u.ID), Op: ragnar.Trace(), Err: err}
	case err != nil:
		log.Fatal(err)
		return err
	default:
	}

	return nil
}

func (db DB) ReadAll(us *[]ragnar.User) error {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return &ragnar.Error{Code: ragnar.ENOTFOUND, Op: ragnar.Trace(), Err: err}
	}

	for rows.Next() {
		u := ragnar.User{}
		err = rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt)
		if err != nil {
			return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: fmt.Sprintf("No user with that ID: %s", u.ID), Op: ragnar.Trace(), Err: err}
		}
		*us = append(*us, u)
	}

	return nil
}

func (db DB) Update(u *ragnar.User) error {
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
func (db DB) Delete(u *ragnar.User) error {
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
		return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
	}

	return nil
}
