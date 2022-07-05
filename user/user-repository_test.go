package user

import (
	"reflect"
	"testing"
)

func createTestUserRepository() (*UserRepository, error) {
	r := NewUserRepository()
	users := []*User{
		{Id: "1", Name: "test1", Password: "test1"},
		{Id: "2", Name: "test2", Password: "test2"},
		{Id: "3", Name: "test3", Password: "test3"},
	}
	for _, user := range users {
		err := r.Add(user)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func TestAdd(t *testing.T) {
	r, err := createTestUserRepository()
	if err != nil {
		t.Error("expected err to be nil")
	}
	registeredUser := r.GetAll()[0]
	err = r.Add(registeredUser)
	if err != ErrUserAlreadyRegistered {
		t.Errorf("expected error to be ErrUserAlreadyRegistered, got: %s", err)
	}
}

func TestById(t *testing.T) {
	r, err := createTestUserRepository()
	if err != nil {
		t.Error("expected err to be nil")
	}
	_, err = r.GetById("1")
	if err != nil {
		t.Error("expected err to be nil")
	}
	_, err = r.GetById("1111")
	if err != ErrUserNotFound {
		t.Errorf("expected error to be ErrUserNotFound, got: %s", err)
	}
}

func TestGetByName(t *testing.T) {
	r, err := createTestUserRepository()
	if err != nil {
		t.Error("expected err to be nil")
	}
	_, err = r.getIndexByName("test1")
	if err != nil {
		t.Error("expected err to be nil")
	}
	_, err = r.getIndexByName("test11111")
	if err != ErrUserNotFound {
		t.Errorf("expected error to be ErrUserNotFound, got: %s", err)
	}
}

func TestDelete(t *testing.T) {
	r, err := createTestUserRepository()
	if err != nil {
		t.Error("expected err to be nil")
	}
	err = r.Delete("1")
	if err != nil {
		t.Error("expected err to be nil")
	}
	usersLength := len(r.GetAll())
	if usersLength != 2 {
		t.Errorf("expected usersLength to be 2, got %d", usersLength)
	}
	err = r.Delete("111111111")
	if err != ErrUserNotFound {
		t.Errorf("expected error to be ErrUserNotFound, got: %s", err)
	}
}

func TestGetAll(t *testing.T) {
	r, err := createTestUserRepository()
	if err != nil {
		t.Error("expected err to be nil")
	}
	if !reflect.DeepEqual(r.GetAll(), r.GetAll()) {
		t.Error("expected users to be the same as repository.users")
	}
}
