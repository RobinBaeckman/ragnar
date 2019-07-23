package mock

import (
	"database/sql"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RobinBaeckman/ragnar"
	"github.com/alicebob/miniredis"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() (DB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	return DB{db, &mock}, nil
}

func NewMemDB() MemDB {
	r, _ := miniredis.Run()

	return MemDB{r}
}

type MemDB struct {
	Redis *Miniredis
}

type DB struct {
	DB   *sql.DB
	Mock *sqlmock.sqlmock
}

func (s DB) Create(u *ragnar.User) error {
	s.Mock.ExpectExec("INSERT INTO users").
		WithArgs(u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Role).
		WillReturnResult(NewResult(1, 1))

	if err := mock.ExpectationsWereMet(); err != nil {
		return &ragnar.Error{Code: ragnar.ECONFLICT, Message: "Username already exists", Op: ragnar.Trace(), Err: err}
	}

	return nil
}

func (s DB) Read(u *ragnar.User) error {
	err := s.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE id=?", u.ID).Scan(&u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: "No user with that ID.", Op: ragnar.Trace(), Err: err}
	case err != nil:
		log.Fatal(err)
	}

	return nil
}

func (s DB) ReadByEmail(u *ragnar.User) error {
	mock.ExpectQuery("SELECT (email, password, first_name, last_name, role) FROM users WHERE id = ?").
		WithArgs(5).
		WillReturnRows(rs)

	err := s.QueryRow("SELECT email, password, first_name, last_name, role FROM users WHERE email=?", u.Email).Scan(&u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role)
	switch {
	case err == sql.ErrNoRows:
		return &ragnar.Error{Code: ragnar.ENOTFOUND, Message: "No user with that ID.", Op: ragnar.Trace(), Err: err}
	case err != nil:
		log.Fatal(err)
		return err
	default:
	}

	return nil
}

func (s DB) ReadAll(us *[]ragnar.User) error {
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

func (s DB) Update(u *ragnar.User) error {
	stmtIns, err := s.Prepare("UPDATE users set email=?, password=?, first_name=?, last_name=?, role=? where id=?")
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
func (s DB) Delete(u *ragnar.User) error {
	stmtIns, err := s.Prepare("DELETE from users where id=?")
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

func (mdb MemDB) Get(k string) (string, error) {
	v, err := mdb.Redis.Get(k).Result()
	if err != nil {
		return v, &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
	}

	return v, err
}

func (mdb MemDB) Set(k string, v interface{}, expiration time.Duration) error {
	if err := mdb.Redis.Set(k, v, 0).Err(); err != nil {
		return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
	}

	return nil
}

func (mdb MemDB) Del(v string) error {
	if err := mdb.Redis.Del(v).Err(); err != nil {
		return &ragnar.Error{Code: ragnar.EINTERNAL, Op: ragnar.Trace(), Err: err}
	}

	return nil
}
