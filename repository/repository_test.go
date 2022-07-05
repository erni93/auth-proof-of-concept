package repository

import "testing"

type TestItem struct {
	Id   string
	Name string
}

func createRepository() *Repository[TestItem] {
	r := NewRepository[TestItem]()
	items := []*TestItem{
		{Id: "1", Name: "test1"},
		{Id: "2", Name: "test2"},
		{Id: "3", Name: "test3"},
	}
	for _, item := range items {
		r.Add(item)
	}
	return r
}

func getById(item *TestItem, value string) bool {
	return item.Id == value
}

func TestAdd(t *testing.T) {
	r := createRepository()
	length := len(r.items)
	if length != 3 {
		t.Errorf("expected length to be 3, got: %d", length)
	}
}

func TestGetAll(t *testing.T) {
	r := createRepository()
	allItems := r.GetAll()
	if len(r.items) != len(allItems) {
		t.Errorf("expected length to be  %d, got: %d", len(allItems), len(r.items))
	}
}

func TestGetItem(t *testing.T) {
	r := createRepository()
	item1, i, err := r.GetItem(getById, "1")
	if item1.Name != "test1" {
		t.Errorf("expected name to be test1, got %s", item1.Name)
	}
	if i != 0 {
		t.Errorf("expected index to be 0, got %d", i)
	}
	if err == ErrItemNotFound {
		t.Error("error not expected, got ErrItemNotFound")
	}
	_, _, err = r.GetItem(getById, "aaaaaa")
	if err != ErrItemNotFound {
		t.Error("expected error ErrItemNotFound, got nil")
	}
}

func TestDeleteItem(t *testing.T) {
	r := createRepository()
	_, i, err := r.GetItem(getById, "1")
	if i != 0 {
		t.Errorf("expected index to be 0, got %d", i)
	}
	if err == ErrItemNotFound {
		t.Error("error not expected, got ErrItemNotFound")
	}
	r.Delete(i)
	_, _, err = r.GetItem(getById, "1")
	if err != ErrItemNotFound {
		t.Error("expected error ErrItemNotFound, got nil")
	}
}
