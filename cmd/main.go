package main

import (
	"log"
	"redis-go/internal/server"
)

func main() {
	srv := server.New(":6379")

	log.Println("redis-go listening on :6379")

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
