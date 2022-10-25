package main

import (
	_ "bookstore/internal/store"
	"bookstore/server"
	"bookstore/store/factory"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	store, err := factory.New("mem")
	if err != nil {
		panic(err)
	}

	srv := server.NewBookStoreServer(":8051", store)
	errChan, err := srv.ListenAndServe()
	if err != nil {
		log.Println("web server start failed:", err)
		return
	}
	log.Println("web server start ok")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-errChan:
		log.Println("web server run failed:", err)
		return
	case <-ch:
		log.Println("bookstore is exiting...")
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		err = srv.Shutdown(ctx)
	}

	if err != nil {
		log.Println("bookstore exit error:", err)
		return
	}
	log.Println("bookstore exit ok")
}
