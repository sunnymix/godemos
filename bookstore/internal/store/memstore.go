package store

import (
	storeDef "bookstore/store"
	factory "bookstore/store/factory"
	"fmt"
	"sync"
)

func init() {
	factory.Register("mem", &MemStore{
		books: make(map[string]*storeDef.Book),
	})
}

type MemStore struct {
	sync.RWMutex
	books map[string]*storeDef.Book
}

func (store *MemStore) Create(book *storeDef.Book) error {
	// FIXME
	store.Lock()
	defer store.Unlock()
	store.books[book.Id] = book
	return nil
}

func (store *MemStore) Update(book *storeDef.Book) error {
	// FIXME
	store.Lock()
	defer store.Unlock()
	store.books[book.Id] = book
	return nil
}

func (store *MemStore) Get(bookId string) (storeDef.Book, error) {
	book := store.books[bookId]
	if book != nil {
		return *book, nil
	}
	return storeDef.Book{}, fmt.Errorf("store: book with id '%s' is not found", bookId)
}

func (store *MemStore) GetAll() ([]storeDef.Book, error) {
	return []storeDef.Book{}, nil
}

func (store *MemStore) Delete(bookId string) error {
	return nil
}
