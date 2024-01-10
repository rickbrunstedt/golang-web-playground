package router

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
)

type Session struct {
	id         string
	authorized bool
}

var sessionMap = map[string]Session{}

func randomText() (string, error) {
	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

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
		return Session{}, err
	}
	return newSession, nil
}

func HandleSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			session, err := createSession()
			if err != nil {
				http.Error(w, "Internal server error", 500)
			}
			cookie := http.Cookie{Name: "session", Value: session.id}
			http.SetCookie(w, &cookie)
		} else {
			_, err = getOrCreateSession(sessionCookie.Value)
			if err != nil {
				http.Error(w, "Internal server error", 500)
			}
		}
		next.ServeHTTP(w, r)
	}
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logged connection from", r.RemoteAddr, r.URL.Path, r.Method, r.URL.Query())
		next.ServeHTTP(w, r)
	}
}
