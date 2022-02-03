package user

import (
	"authGo/user"
	"errors"
)

var ErrUserNotFound = errors.New("user repository: user not found")
var ErrUserAlreadyRegistered = errors.New("user repository: user already registered")

type UserRepository struct {
	users []*user.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: make([]*user.User, 0)}
}

func (r *UserRepository) Add(user *user.User) error {
	i, _ := r.getIndexByName(user.Name)
	if i != -1 {
		return ErrUserAlreadyRegistered
	}
	r.users = append(r.users, user)
	return nil
}

func (r *UserRepository) GetById(id string) (*user.User, error) {
	i, err := r.getIndexById(id)
	if err != nil {
		return nil, err
	}
	return r.users[i], ErrUserNotFound
}

func (r *UserRepository) GetByName(name string) (*user.User, error) {
	i, err := r.getIndexByName(name)
	if err != nil {
		return nil, err
	}
	return r.users[i], ErrUserNotFound
}

func (r *UserRepository) Delete(id string) error {
	i, err := r.getIndexById(id)
	if err != nil {
		return err
	}
	lastIndex := len(r.users) - 1
	r.users[i] = r.users[lastIndex]
	r.users[lastIndex] = nil
	r.users = r.users[:lastIndex]
	return nil
}

func (r *UserRepository) getIndexById(id string) (int, error) {
	for i, user := range r.users {
		if user.Id == id {
			return i, nil
		}
	}
	return -1, ErrUserNotFound
}

func (r *UserRepository) getIndexByName(name string) (int, error) {
	for i, user := range r.users {
		if user.Name == name {
			return i, nil
		}
	}
	return -1, ErrUserNotFound
}
