package repository

type ComparableFunc[T any] func(item *T, value string) bool

type Repository[T any] struct {
	items []*T
}

func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{items: make([]*T, 0)}
}

func (r *Repository[T]) GetItem(comparable ComparableFunc[T], value string) (*T, int) {
	for i, item := range r.items {
		if comparable(item, value) {
			return item, i
		}
	}
	return nil, -1
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
	r.items = r.items[:lastIndex]
}
