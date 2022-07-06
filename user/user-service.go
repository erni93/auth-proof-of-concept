package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository        *UserRepository
	passwordValidator PasswordValidator
}

type PasswordValidator interface {
	compareHashAndPassword(hashedPassword, password []byte) error
	generateFromPassword(password []byte, cost int) ([]byte, error)
}

type ServicePasswordValidator struct{}

func (spv ServicePasswordValidator) compareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (spv ServicePasswordValidator) generateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func NewUserService() *UserService {
	return &UserService{repository: NewUserRepository(), passwordValidator: ServicePasswordValidator{}}
}

func (s *UserService) IsPasswordValid(name string, password string) bool {
	user, err := s.repository.GetByName(name)
	if err != nil {
		return false
	}
	if err = s.passwordValidator.compareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func (s *UserService) CreateUser(name string, password string) error {
	id := uuid.New()
	passwordHash, err := s.getPasswordHash(password)
	if err != nil {
		return err
	}
	err = s.repository.Add(&User{Id: id.String(), Name: name, Password: passwordHash})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(id string) error {
	return s.repository.Delete(id)
}

func (s *UserService) GetAllUsers() []*User {
	return s.repository.GetAll()
}

func (s *UserService) getPasswordHash(password string) (string, error) {
	passwordBytes, err := s.passwordValidator.generateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil
}
