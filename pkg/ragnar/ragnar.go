package ragnar

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
	return fmt.Sprintf("%s:%d:%s", file, line, f.Name())
}

type User struct {
	ID           string
	Email        string
	Password     string
	PasswordHash []byte `json:"-"`
	FirstName    string
	LastName     string
	Role         string
	CreatedAt    string `json:"-"`
}

type Error struct {
	// Machine-readable error code.
	Code string

	// Human-readable message.
	Message string

	// The operation being performed, usually the method
	// being invoked (Get, Put, etc.)..
	Op string

	// Logical operation and nested error.
	Err error
}

// Error returns the string representation of the error message.
// TODO: create a formated error message
func (e *Error) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s, Op: %s, Err: %s", e.Code, e.Message, e.Op, e.Err)
}

// Application error codes.
const (
	ECONFLICT     = "conflict"    // action cannot be performed
	EINTERNAL     = "internal"    // internal error
	EINVALID      = "invalid"     // validation failed
	EFORBIDDEN    = "not_allowed" // authentication failed
	EUNAUTHORIZED = "not_authorized"
	ENOTFOUND     = "not_found" // entity does not exist
)

// Application user messages
const (
	EINTERNAL_MSG = "Oops something went wrong"
)

type DB interface {
	Create(*User) error
	Read(*User) error
	ReadByEmail(*User) error
	ReadAll(*[]User) error
	Update(*User) error
	Delete(*User) error
	CleanupTables() error
	Close()
}

type MemDB interface {
	Set(string, interface{}, time.Duration) error
	Del(string) error
	Get(string) (string, error)
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
}
