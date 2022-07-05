package repository

import (
	"strings"
	"testing"
)

type TestItem struct {
	Id   string
	Name string
}

func createRepository() *Repository[TestItem] {
	r := NewRepository[TestItem]()
	items := []*TestItem{
		{Id: "1", Name: "test1"},
		{Id: "2", Name: "test2"},
		{Id: "3", Name: "super test3"},
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

	t.Run("get different items at different positions", func(t *testing.T) {
		items := []struct {
			item  *TestItem
			index int
		}{
			{r.items[0], 0},
			{r.items[1], 1},
			{r.items[2], 2},
		}

		for _, testItem := range items {
			got, i := r.GetItem(getById, testItem.item.Id)
			if i != testItem.index {
				t.Errorf("expected %d to be %d", i, testItem.index)
			}
			if got != testItem.item {
				t.Errorf("expected %v to be %v", got, testItem.item)
			}
		}
	})

	t.Run("item not found", func(t *testing.T) {
		_, i := r.GetItem(getById, "aaaaaa")
		if i != -1 {
			t.Errorf("expected index to be -1, got %d", i)
		}
	})

	t.Run("allow custom comparable functions", func(t *testing.T) {
		nameContains := func(item *TestItem, value string) bool {
			return strings.Contains(item.Name, value)
		}
		item, _ := r.GetItem(nameContains, "super")
		if item == nil {
			t.Error("expected item to not be nil")
		}
	})
}

func TestDeleteItem(t *testing.T) {
	r := createRepository()
	_, i := r.GetItem(getById, "1")
	if i != 0 {
		t.Errorf("expected index to be 0, got %d", i)
	}
	r.Delete(i)
	_, i = r.GetItem(getById, "1")
	if i != -1 {
		t.Errorf("expected index to be -1, got %d", i)
	}
}
