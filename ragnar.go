package ragnar

// Product is a product that is being sold
type User struct {
	ID        string
	Email     string
	Password  []byte
	FirstName string
	LastName  string
	Role      string
}

type UserService interface {
	Store(*User) error
	Get(*User) error
	GetByEmail(string) (*User, error)
	GetAll(*[]User) error
}

type In struct {
	ID        string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      string
}
