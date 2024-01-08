package main

import (
	"app/web/router"

	"fmt"
	"log"
	"net/http"
)

var PORT = ":8080"

func main() {
	mux := http.NewServeMux()
	middlewares := []func(http.HandlerFunc) http.HandlerFunc{router.HandleSession}
	mux.HandleFunc("/", router.ApplyMiddlewares(middlewares, routes))

	fmt.Println("Starting server on port", PORT)
	err := http.ListenAndServe(PORT, mux)
	log.Fatal(err)
}

func routes(w http.ResponseWriter, r *http.Request) {
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
	case "/authorize":
		authorize(w, r)
	case "/check-auth":
		checkAuth(w, r)
	default:
		http.NotFound(w, r)
	}
}

func checkAuth(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Cookie")
	w.Header().Add("Cookie", "test=1234")
	w.Header().Set("Authorization", "3214")

	if auth == "1234" {
		fmt.Fprintf(w, "Authorized")
	} else {
		fmt.Fprintf(w, "Not Authorized")
	}
}

func authorize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Authorization", "1234")
	http.Redirect(w, r, "/check-auth", 308)
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
