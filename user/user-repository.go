package user

import (
	"authGo/repository"
	"errors"
)

var ErrUserNotFound = errors.New("user repository: user not found")
var ErrUserAlreadyRegistered = errors.New("user repository: user already registered")

func getById(user *User, value string) bool {
	return user.Id == value
}

func getByName(user *User, value string) bool {
	return user.Name == value
}

type UserRepository struct {
	repository *repository.Repository[User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{repository: repository.NewRepository[User]()}
}

func (r *UserRepository) Add(user *User) error {
	i, _ := r.getIndexByName(user.Name)
	if i != -1 {
		return ErrUserAlreadyRegistered
	}
	r.repository.Add(user)
	return nil
}

func (r *UserRepository) GetById(id string) (*User, error) {
	user, _, err := r.getUser(getById, id)
	return user, err
}

func (r *UserRepository) GetByName(name string) (*User, error) {
	user, _, err := r.getUser(getByName, name)
	return user, err
}

func (r *UserRepository) Delete(id string) error {
	i, err := r.getIndexById(id)
	if err != nil {
		return err
	}
	r.repository.Delete(i)
	return nil
}

func (r *UserRepository) GetAll() []*User {
	return r.repository.GetAll()
}

func (r *UserRepository) getUser(comparableFunc repository.ComparableFunc[User], value string) (*User, int, error) {
	user, i := r.repository.GetItem(comparableFunc, value)
	if i == -1 {
		return nil, i, ErrUserNotFound
	}
	return user, i, nil
}

func (r *UserRepository) getIndexById(id string) (int, error) {
	_, i, err := r.getUser(getById, id)
	return i, err
}

func (r *UserRepository) getIndexByName(name string) (int, error) {
	_, i, err := r.getUser(getByName, name)
	return i, err
}
