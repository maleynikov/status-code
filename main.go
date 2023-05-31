package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "8080"
	}

	if err := runServer(port); err != nil {
		log.Fatal("Server cannot start")
		os.Exit(1)
	}
}

func handlerOK(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handlerBadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
}

func handlerServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}

func contentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func runServer(port string) (err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ok", handlerOK)
	mux.HandleFunc("/bad-request", handlerBadRequest)
	mux.HandleFunc("/internal-server-error", handlerServerError)
	handler := contentType(mux)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("mux server is starting at :%s\n", port)
	err = server.ListenAndServe()

	return nil
}
