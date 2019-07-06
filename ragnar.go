package ragnar

type User struct {
	ID           string
	Email        string
	Password     string
	PasswordHash []byte `json:"-"`
	FirstName    string
	LastName     string
	Role         string
}

type Message struct {
	ID   string
	Body string
}

type UserService interface {
	Create(*User) error
	Read(*User) error
	ReadByEmail(*User) error
	ReadAll(*[]User) error
	Update(*User) error
	Delete(*User) error
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
