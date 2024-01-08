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
	fmt.Println("Starting server on port", PORT)
	err := http.ListenAndServe(PORT, mux)
	log.Fatal(err)
}

func applyMiddlewares(middlewares []func(http.HandlerFunc) http.HandlerFunc, next http.HandlerFunc) http.HandlerFunc {
	if len(middlewares) < 1 {
		return next
	}

	return middlewares[0](applyMiddlewares(middlewares[1:], next))
}

func randomText() (string, error) {
	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

type Session struct {
	id         string
	authorized bool
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

func getOrCreateSession(sessionId string) (Session, error) {
	session, ok := sessionMap[sessionId]
	if ok {
		return session, nil
	}
	newSession, err := createSession()
	if err != nil {
		log.Println("Error creating session:", err)
		return Session{}, err
	}
	return newSession, nil
}

func handleSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			session, err := createSession()
			if err != nil {
				log.Println("Error creating session:", err)
				http.Error(w, "Internal server error", 500)
				return
			}
			cookie := http.Cookie{Name: "session", Value: session.id}
			http.SetCookie(w, &cookie)
		}
		session, err := getOrCreateSession(sessionCookie.Value)
		if err != nil {
			log.Println("Error getting session:", err)
			http.Error(w, "Internal server error", 500)
			return
		}
		log.Println("Session:", session)
		next.ServeHTTP(w, r)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Just a dummy thing for now
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logged connection from", r.RemoteAddr, r.URL.Path, r.Method, r.URL.Query())
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
