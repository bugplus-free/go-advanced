package main

import (
	"log"
	"net/http"
	"os"

	"go-advanced/day1/homepage"
	"go-advanced/day1/server"
)

var (
	GcukServiceAddr = os.Getenv("GCUK_SERVICE_ADDR")
)

func main() {
	logger := log.New(os.Stdout, "gcuk", log.LstdFlags|log.Lshortfile)
	h := homepage.NewHandlers(logger)

	mux := http.NewServeMux()
	h.SetupRoutes(mux)

	srv := server.New(mux, "127.0.0.1:8080")

	logger.Println("server starting!!!")
	// err := http.ListenAndServe(":8080", mux)
	// err := srv.ListenAndServe()
	err := srv.ListenAndServeTLS("key/req.pem", "key/decrypted_key.pem")
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}

}
