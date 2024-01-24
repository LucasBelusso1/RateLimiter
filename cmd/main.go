package main

import (
	"log"
	"net/http"

	"github.com/LucasBelusso1/go-ratelimiter/pkg/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	err := http.ListenAndServe(":8080", middleware.AccessLimitMiddleware(mux))

	log.Println("Listening on port :8080")
	if err != nil {
		log.Fatalf("Could not listen on :8080: %v\n", err)
	}
}
