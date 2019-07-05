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

type UserService interface {
	Create(*User) error
	Read(*User) error
	ReadByEmail(*User) error
	ReadAll(*[]User) error
	Update(*User) error
	Delete(*User) error
}
