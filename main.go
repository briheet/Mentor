package main

import (
	"log"
)

func main() {
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = "3000" // Default to port 3000 if PORT is not set
	//}

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}

