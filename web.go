package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
)

var PORT = ":8080"

func main() {
	mux := http.NewServeMux()
	middlewares := []func(http.HandlerFunc) http.HandlerFunc{handleSession}
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

var randomReader = rand.Reader

func randomText() (string, error) {
	buf := make([]byte, 32)
	_, err := io.ReadFull(randomReader, buf)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

type Session struct {
	id         string
	authorized bool
	// name string
}

var sessionMap = map[string]Session{}

func createSession() (Session, error) {
	nextSessionId, err := randomText()

	if err != nil {
		return Session{}, err
	}

	session := Session{id: nextSessionId, authorized: false}
	sessionMap[nextSessionId] = session
	return session, nil
}

func handleSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")

		if err != nil {
			session, err := createSession()

			if err != nil {
				log.Printf("Error creating session: %s", err)
				http.Error(w, "Internal server error", 500)
				return
			}
			cookie := http.Cookie{Name: "session", Value: session.id}
			http.SetCookie(w, &cookie)
		} else {
			log.Printf("Session cookie found: %s", sessionCookie.Name)
			log.Printf("Session cookie found: %s", sessionCookie.Value)
			session, ok := sessionMap[sessionCookie.Value]

			if !ok {
				log.Printf("Session not found: %s", sessionCookie.Value)
				http.Error(w, "Internal server error", 500)
				return
			}

			if session.authorized {
				log.Printf("Session authorized: %s", sessionCookie.Value)
				next.ServeHTTP(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Just a dummy thing for now
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Authenticated connection from %s %s", r.RemoteAddr, r.URL.Path)
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
	// log.Printf("Auth: %s", auth)
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
