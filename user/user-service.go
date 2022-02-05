package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *UserRepository
}

func NewUserService() *UserService {
	return &UserService{repository: NewUserRepository()}
}

func (s *UserService) CreateUser(name string, password string) error {
	id := uuid.New()
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.repository.Add(&User{Id: id.String(), Name: name, Password: string(passwordByte)})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) IsPasswordValid(name string, password string) bool {
	user, err := s.repository.GetByName(name)
	if err != nil {
		return false
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func (s *UserService) DeleteUser(id string) error {
	return s.repository.Delete(id)
}

func (s *UserService) GetAllUsers() []*User {
	return s.repository.GetAll()
}
