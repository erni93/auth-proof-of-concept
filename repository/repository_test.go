package repository

import (
	"errors"
	"fmt"
	"testing"
)

type TestItem struct {
	Id   string
	Name string
}

func (item TestItem) String() string {
	return fmt.Sprintf("(id: %s,  name:  %s)", item.Id, item.Name)
}

func createTestItems() []TestItem {
	return []TestItem{
		{Id: "1", Name: "Foo"},
		{Id: "2", Name: "Bar"},
	}
}

func TestInit(t *testing.T) {
	if r := NewRepository(); r == nil {
		t.Errorf("expected NewRepository to return Repository, got %s", r)
	}
}

func TestAdd(t *testing.T) {
	t.Run("add value with Id property", func(t *testing.T) {
		r := NewRepository()
		err := r.Add(createTestItems()[0])
		if err != nil {
			t.Errorf("expected to be successful, got: %s", err)
		}
	})
	t.Run("get error when adding a value that is not struct", func(t *testing.T) {
		r := NewRepository()
		err := r.Add("random value :)")
		if !errors.Is(err, ErrItemNotStruct) {
			t.Errorf("expected error to be ErrItemNotStruct, got: %s", err)
		}
	})
	t.Run("get err adding value without Id property", func(t *testing.T) {
		r := NewRepository()
		item := struct{ name string }{name: "custom struct"}
		err := r.Add(item)
		if !errors.Is(err, ErrItemWithoutId) {
			t.Errorf("expected error to be ErrItemWithoutId, got: %s", err)
		}
	})
}

func TestGet(t *testing.T) {
	r := NewRepository()
	err := r.Add(createTestItems()[0])
	if err != nil {
		t.Errorf("expected err to be nil, got: %s", err)
	}
	t.Run("get item by id", func(t *testing.T) {
		item, _ := r.Get("1")
		got := item.(TestItem)
		want := createTestItems()[0]
		if got != want {
			t.Errorf("Expected %s to be %s", got, want)
		}
	})
	t.Run("get error if id was not found", func(t *testing.T) {
		_, err := r.Get("random id :)")
		if !errors.Is(err, ErrItemNotFound) {
			t.Errorf("expected error to be ErrItemNotFound, got: %s", err)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("delete item by id", func(t *testing.T) {
		r := NewRepository()
		err := r.Add(createTestItems()[0])
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}

		err = r.Delete("1")
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}

		got := len(r.items)
		want := 0
		if got != want {
			t.Errorf("Expected %d to be %d", got, want)
		}
	})
	t.Run("get error if id was not found", func(t *testing.T) {
		r := NewRepository()
		err := r.Add(createTestItems()[0])
		if err != nil {
			t.Errorf("expected err to be nil, got: %s", err)
		}

		err = r.Delete("random id :)")
		if !errors.Is(err, ErrItemNotFound) {
			t.Errorf("expected error to be ErrItemNotFound, got: %s", err)
		}
	})
}
