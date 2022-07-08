package user

import (
	"errors"
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
		{Id: "1", Name: "test1", Password: "test1", IsAdmin: true},
		{Id: "2", Name: "test2", Password: "test2", IsAdmin: false},
		{Id: "3", Name: "test3", Password: "test3", IsAdmin: false},
	}
	service := NewUserService()
	for _, user := range users {
		err := service.CreateUser(user.Name, user.Password, user.IsAdmin)
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
		err := s.CreateUser("test1", "test1", true)
		if err != ErrUserAlreadyRegistered {
			t.Errorf("expected error to be ErrUserAlreadyRegistered, got: %s", err)
		}
	})
	t.Run("Assign isAdmin to the new user", func(t *testing.T) {
		err := s.CreateUser("admin", "admin", true)
		if err != nil {
			t.Errorf("expected error to be nil, got: %s", err)
		}
		user, err := s.GetRepository().GetByName("admin")
		if err != nil {
			t.Errorf("expected error to be nil, got: %s", err)
		}
		if user.IsAdmin != true {
			t.Error("expected IsAdmin to be true, got")
		}
	})
	t.Run("Error hashing user password", func(t *testing.T) {
		hashError := errors.New("hash error :( ")
		mockPasswordValidator := MockPasswordValidator{
			ErrorToReturn: hashError,
		}
		s.passwordValidator = mockPasswordValidator
		err := s.CreateUser("test4", "test4", true)
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

func TestGetRepository(t *testing.T) {
	s, err := createTestUserService()
	if err != nil {
		t.Errorf("error creating UserService, %s", err)
	}

	repository := s.GetRepository()
	if repository != s.repository {
		t.Error("expected repository to be the same as service repository")
	}
}
