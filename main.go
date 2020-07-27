package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	} else {
		log.Printf("Server listening to port :5000")
	}
}
