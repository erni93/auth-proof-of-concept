package user

import (
	"errors"
	"reflect"
	"testing"
)

type Login struct {
	Name     string
	Password string
}

type MockPasswordValidator struct {
	ErrorToReturn error
}

func (spv MockPasswordValidator) compareHashAndPassword(hashedPassword, password []byte) error {
	return spv.ErrorToReturn
}

func (spv MockPasswordValidator) generateFromPassword(password []byte, cost int) ([]byte, error) {
	return nil, spv.ErrorToReturn
}

func createTestUserService() (*UserService, error) {
	users := []*User{
		{Id: "1", Name: "test1", Password: "test1"},
		{Id: "2", Name: "test2", Password: "test2"},
		{Id: "3", Name: "test3", Password: "test3"},
	}
	service := NewUserService()
	for _, user := range users {
		err := service.CreateUser(user.Name, user.Password)
		if err != nil {
			return nil, err
		}
	}
	return service, nil

}

func TestCreateUser(t *testing.T) {
	s, err := createTestUserService()
	if err != nil {
		t.Errorf("error creating UserService, %s", err)
	}
	t.Run("Create user with the same name", func(t *testing.T) {
		err := s.CreateUser("test1", "test1")
		if err != ErrUserAlreadyRegistered {
			t.Errorf("expected error to be ErrUserAlreadyRegistered, got: %s", err)
		}
	})
	t.Run("Error hashing user password", func(t *testing.T) {
		hashError := errors.New("hash error :( ")
		mockPasswordValidator := MockPasswordValidator{
			ErrorToReturn: hashError,
		}
		s.passwordValidator = mockPasswordValidator
		err := s.CreateUser("test4", "test4")
		if err != hashError {
			t.Errorf("expected error to be %s", hashError)
		}
		isPasswordValid := s.IsPasswordValid("test4", "test4")
		if isPasswordValid != false {
			t.Error("expected isPasswordValid to be false")
		}
	})
}

func TestIsPasswordValid(t *testing.T) {
	s, err := createTestUserService()
	if err != nil {
		t.Errorf("error creating UserService, %s", err)
	}
	loginTests := []struct {
		login Login
		want  bool
	}{
		{Login{Name: "", Password: ""}, false},
		{Login{Name: "test1", Password: "abc"}, false},
		{Login{Name: "test2", Password: "test1"}, false},
		{Login{Name: "test1", Password: "test11"}, false},
		{Login{Name: "test1", Password: "test11"}, false},
		{Login{Name: "test1", Password: "tEsT1"}, false},
		{Login{Name: "test1", Password: "test1"}, true},
		{Login{Name: "test2", Password: "test2"}, true},
		{Login{Name: "test3", Password: "test3"}, true},
	}

	for _, loginTest := range loginTests {
		got := s.IsPasswordValid(loginTest.login.Name, loginTest.login.Password)
		if got != loginTest.want {
			t.Errorf("login %s password %s, got %t want %t", loginTest.login.Name, loginTest.login.Password, got, loginTest.want)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	s, err := createTestUserService()
	if err != nil {
		t.Errorf("error creating UserService, %s", err)
	}

	id := s.repository.GetAll()[0].Id
	err = s.DeleteUser(id)
	if err != nil {
		t.Error("expected err to be nil")
	}
	usersLength := len(s.repository.GetAll())
	if usersLength != 2 {
		t.Errorf("expected usersLength to be 2, got %d", usersLength)
	}
	err = s.DeleteUser("123")
	if err == nil {
		t.Error("expected err to not be nil")
	}
	usersLength = len(s.repository.GetAll())
	if usersLength != 2 {
		t.Errorf("expected usersLength to be 2, got %d", usersLength)
	}
}

func TestGetAllUsers(t *testing.T) {
	s, err := createTestUserService()
	if err != nil {
		t.Errorf("error creating UserService, %s", err)
	}

	users := s.GetAllUsers()
	if !reflect.DeepEqual(users, s.repository.GetAll()) {
		t.Error("expected users to be the same as repository.users")
	}
}
