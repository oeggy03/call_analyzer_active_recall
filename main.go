package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/daily", dailyHandler)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func newRandom() *rand.Rand {
	seed := uint64(time.Now().UnixNano())
	return rand.New(rand.NewPCG(seed, seed^uint64(os.Getpid())))
}
