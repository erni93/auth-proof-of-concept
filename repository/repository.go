package repository

import "errors"

var ErrItemNotFound = errors.New("repository: item not found")

type ComparableFunc[T any] func(item *T, value string) bool

type Repository[T any] struct {
	items []*T
}

func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{items: make([]*T, 0)}
}

func (r *Repository[T]) GetItem(comparable ComparableFunc[T], value string) (*T, int, error) {
	for i, item := range r.items {
		if comparable(item, value) {
			return item, i, nil
		}
	}
	return nil, -1, ErrItemNotFound
}

func (r *Repository[T]) GetAll() []*T {
	return r.items
}

func (r *Repository[T]) Add(item *T) {
	r.items = append(r.items, item)
}

func (r *Repository[T]) Delete(index int) {
	lastIndex := len(r.items) - 1
	r.items[index] = r.items[lastIndex]
	r.items = r.items[:lastIndex-1]
}
