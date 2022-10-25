package server

import (
	"bookstore/server/middleware"
	"bookstore/store"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type BookStoreServer struct {
	store store.Store
	srv   *http.Server
}

func NewBookStoreServer(addr string, store store.Store) *BookStoreServer {
	srv := &BookStoreServer{
		store: store,
		srv: &http.Server{
			Addr: addr,
		},
	}
	router := mux.NewRouter()
	router.HandleFunc("/book", srv.createBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", srv.getBookHandler).Methods("GET")
	srv.srv.Handler = middleware.Logging(middleware.Validating(router))
	return srv
}

func (bss *BookStoreServer) Shutdown(ctx context.Context) error {
	return bss.srv.Shutdown(ctx)
}

func (bss *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		err = bss.srv.ListenAndServe()
		errChan <- err
	}()
	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}

func (bss *BookStoreServer) createBookHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bss.store.Create(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bss *BookStoreServer) getBookHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "no id found in request", http.StatusBadRequest)
		return
	}
	book, err := bss.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response(w, book)
}

func response(w http.ResponseWriter, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
