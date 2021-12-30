package repository

import (
	"errors"
	"reflect"
)

type Repository interface {
	Init() interface{}
	Add(value interface{}) interface{}
	Delete(id string) error
	Get(id string) error
}

type RepositoryBase struct {
	items []interface{}
}

func (r RepositoryBase) Init() *RepositoryBase {
	return &RepositoryBase{items: make([]interface{}, 0)}
}

func (r *RepositoryBase) Add(value *interface{}) {
	r.items = append(r.items, *value)
}

func (r *RepositoryBase) Get(id string) (interface{}, error) {
	i, err := r.getIndex(id)
	if err != nil {
		return nil, err
	}
	return r.items[i], nil
}

func (r *RepositoryBase) Delete(id string) error {
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

func (r *RepositoryBase) getIndex(id string) (int, error) {
	for i, item := range r.items {
		s := reflect.ValueOf(&item).Elem()
		if s.Kind() == reflect.Struct {
			if f := s.FieldByName("Id"); f.IsValid() {
				if f.String() == id {
					return i, nil
				}
				continue
			}
			return -1, errors.New("item doesn't have property 'Id'")
		}
		return -1, errors.New("item type is not 'Struct'")
	}
	return -1, errors.New("item not found")
}
