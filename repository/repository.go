package repository

import (
	"errors"
	"reflect"
)

var ErrItemNotFound = errors.New("repository: item not found")
var ErrItemNotStruct = errors.New("repository: value type is not 'Struct'")
var ErrItemWithoutId = errors.New("repository: value doesn't have property 'Id'")

type Repository struct {
	items []interface{}
}

func NewRepository() *Repository {
	return &Repository{items: make([]interface{}, 0)}
}

func (r *Repository) Add(value interface{}) error {
	if err := r.IsValid(value); err != nil {
		return err
	}
	r.items = append(r.items, value)
	return nil
}

func (r *Repository) Get(id string) (interface{}, error) {
	i, err := r.getIndex(id)
	if err != nil {
		return nil, err
	}
	return r.items[i], nil
}

func (r *Repository) Delete(id string) error {
	i, err := r.getIndex(id)
	if err != nil {
		return err
	}
	lastIndex := len(r.items) - 1
	r.items[i] = r.items[lastIndex]
	r.items[lastIndex] = nil
	r.items = r.items[:lastIndex]
	return nil
}

// We want to prevent a panic error with getIndex trying to get property 'Id'
func (r *Repository) IsValid(item interface{}) error {
	s := reflect.ValueOf(&item).Elem().Elem()
	if s.Kind() != reflect.Struct {
		return ErrItemNotStruct
	}
	if f := s.FieldByName("Id"); !f.IsValid() {
		return ErrItemWithoutId
	}
	return nil
}

func (r *Repository) getIndex(id string) (int, error) {
	for i, item := range r.items {
		if err := r.IsValid(item); err == nil {
			itemId := reflect.ValueOf(&item).Elem().Elem().FieldByName("Id").String()
			if itemId == id {
				return i, nil
			}
		}
	}
	return -1, ErrItemNotFound
}
