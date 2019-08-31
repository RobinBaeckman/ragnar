package rolf

import (
	"fmt"
	"runtime"
	"time"
)

func Trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d:%s\t", file, line, f.Name())
}

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	PasswordHash []byte `json:"-"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Role         string `json:"role"`
	Token        string `json:"token"`
	CreatedAt    string `json:"-"`
}

type Error struct {
	// Machine-readable error code.
	Code int

	// Human-readable message.
	Msg string

	// Sensitive error related data for the operator.
	Op string
}

// Error returns the string representation of the error message.
// TODO: create a formated error message
func (e *Error) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s, Op: %s", e.Code, e.Msg, e.Op)
}

// Application user messages
const (
	EINTERNAL_MSG = "Oops something went wrong"
)

type DB interface {
	Create(*User) error
	Read(*User) error
	ReadAny(*User) error
	ReadByEmail(*User) error
	ReadAll(*[]User) error
	Update(*User) error
	Delete(*User) error
	UpdatePassword(*User) error
	StoreToken(*User) error
	CleanupTables() error
	Close()
}

type MemDB interface {
	SetUser(string, *User, time.Duration) error
	Del(string) error
	GetUser(string) (*User, error)
}

type Mailer interface {
	Send(to string, from string, msg string) error
}

var Env = map[string]string{
	"LOG_PREFIX":  "",
	"HOST":        "",
	"PORT":        "",
	"MYSQL_HOST":  "",
	"MYSQL_USER":  "",
	"MYSQL_PASS":  "",
	"MYSQL_DB":    "",
	"REDIS_HOST":  "",
	"REDIS_PORT":  "",
	"COOKIE_NAME": "",
	"JWT_KEY":     "",
	"PROTO":       "",
	"SMTP_HOST":   "",
	"SMTP_PORT":   "",
	"EMAIL":       "",
}
