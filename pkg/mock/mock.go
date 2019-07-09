package mock

import "github.com/RobinBaeckman/ragnar"

// UserService represents a mock implementation of atla.AdminService.
type UserService struct {
	UserGetFn      func(a *ragnar.User) error
	UserGetInvoked bool

	UserGetAllFn      func(as *[]ragnar.User) error
	UserGetAllInvoked bool

	UserPersistFn      func(a *ragnar.User) error
	UserPersistInvoked bool
}

// User invokes the mock implementation and marks the function as invoked.
func (s *UserService) Get(a *ragnar.User) error {
	s.UserGetInvoked = true
	return s.UserGetFn(a)
}

func (s *UserService) GetAll(as *[]ragnar.User) error {
	s.UserGetAllInvoked = true
	return s.UserGetAllFn(as)
}

func (s *UserService) Persist(a *ragnar.User) error {
	s.UserPersistInvoked = true
	return s.UserPersistFn(a)
}
