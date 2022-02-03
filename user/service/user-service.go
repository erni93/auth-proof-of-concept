package user

import (
	"authGo/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	repository "authGo/user/repository"
)

type UserService struct {
	repository *repository.UserRepository
}

func New() *UserService {
	return &UserService{repository: repository.NewUserRepository()}
}

func (s *UserService) CreateUser(name string, password string) error {
	id := uuid.New()
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.repository.Add(&user.User{Id: id.String(), Name: name, Password: string(passwordByte)})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) isPasswordValid(name string, password string) bool {
	user, err := s.repository.GetByName(name)
	if err != nil {
		return false
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}
