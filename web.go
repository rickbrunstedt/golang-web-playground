package main

import (
	"fmt"
	"log"
	"net/http"
)

var PORT = ":8080"

func main() {
	mux := http.NewServeMux()
	middlewares := []func(http.HandlerFunc) http.HandlerFunc{loggerMiddleware, authMiddleware}
	mux.HandleFunc("/", applyMiddlewares(middlewares, router))

	fmt.Println("Starting server on port " + PORT)
	err := http.ListenAndServe(PORT, mux)
	log.Fatal(err)
}

func applyMiddlewares(middlewares []func(http.HandlerFunc) http.HandlerFunc, next http.HandlerFunc) http.HandlerFunc {
	if len(middlewares) < 1 {
		return next
	}

	return middlewares[0](applyMiddlewares(middlewares[1:], next))
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Just a dummy thing for now
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Authenticated connection from %s %s", r.RemoteAddr, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Logged connection from %s %s %s %s", r.RemoteAddr, r.URL.Path, r.Method, r.URL.Query())
		next.ServeHTTP(w, r)
	}
}

func router(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")

	switch r.URL.Path {
	case "/":
		index(w, r)
	case "/about":
		about(w, r)
	case "/data":
		jsonResponse(w, r)
	case "/query":
		handleQuery(w, r)
	default:
		http.NotFound(w, r)
	}
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Fprintf(w, "Query: %s\n", query)
}

func jsonResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json := `{"status": "ok"}`
	fmt.Fprintf(w, json)
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The about pages")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world! This is the index page.")
}
