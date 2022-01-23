package user

import (
	"authGo/repository"
	"authGo/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repository.Repository
}

func New() *UserService {
	return &UserService{repository: repository.NewRepository()}
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

/* TODO: Complete isPasswordValid
func (s *UserService) isPasswordValid(name string, password string) bool {
}
*/
