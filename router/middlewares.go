package router

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
)

func ApplyMiddlewares(middlewares []func(http.HandlerFunc) http.HandlerFunc, next http.HandlerFunc) http.HandlerFunc {
	if len(middlewares) < 1 {
		return next
	}
	return middlewares[0](ApplyMiddlewares(middlewares[1:], next))
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

func CreateSession() (Session, error) {
	nextSessionId, err := randomText()
	if err != nil {
		return Session{}, err
	}
	session := Session{id: nextSessionId, authorized: false}
	sessionMap[nextSessionId] = session
	return session, nil
}

func GetOrCreateSession(sessionId string) (Session, error) {
	session, ok := sessionMap[sessionId]
	if ok {
		return session, nil
	}
	newSession, err := CreateSession()
	if err != nil {
		log.Println("Error creating session:", err)
		return Session{}, err
	}
	return newSession, nil
}

func HandleSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			session, err := CreateSession()
			if err != nil {
				log.Println("Error creating session:", err)
				http.Error(w, "Internal server error", 500)
				return
			}
			cookie := http.Cookie{Name: "session", Value: session.id}
			http.SetCookie(w, &cookie)
		}
		session, err := GetOrCreateSession(sessionCookie.Value)
		if err != nil {
			log.Println("Error getting session:", err)
			http.Error(w, "Internal server error", 500)
			return
		}
		log.Println("Session:", session)
		next.ServeHTTP(w, r)
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Just a dummy thing for now
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logged connection from", r.RemoteAddr, r.URL.Path, r.Method, r.URL.Query())
		next.ServeHTTP(w, r)
	}
}
