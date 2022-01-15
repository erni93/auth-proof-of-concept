package user

import (
	"authGo/repository"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repository.Repository
}

func NewUserService() *UserService {
	service := &UserService{repository: repository.NewRepository()}
	service.addAdminUser()
	return service
}

func (s *UserService) addAdminUser() {
	id := uuid.New()
	password, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error creating admin password, %s", err)
	}
	err = s.repository.Add(&User{Id: id.String(), Name: "admin", Password: string(password)})
	if err != nil {
		log.Printf("error adding admin user, %s", err)
	}
}
