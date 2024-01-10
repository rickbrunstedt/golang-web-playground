package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []func(http.HandlerFunc) http.HandlerFunc
	routes      map[string]func(http.ResponseWriter, *http.Request)
	Headers     Headers
}

type Headers struct {
	ContentType               string
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
}

func NewRouter() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		middlewares: []func(http.HandlerFunc) http.HandlerFunc{},
		routes:      map[string]func(http.ResponseWriter, *http.Request){},
		Headers:     Headers{},
	}
}

func (r *Router) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	r.routes[path] = handler
}

func (r *Router) Use(middleware func(http.HandlerFunc) http.HandlerFunc) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) Json(w http.ResponseWriter, data interface{}) {
	res, err := formatJson(data)
	if err != nil {
		http.Error(w, "Internal server error", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, res)
}

func (r *Router) Start(port string) error {
	for path, handler := range r.routes {
		withMiddlewares := applyMiddlewares(r.middlewares, handler)
		withHeaders := setDefaultHeaders(r, withMiddlewares)
		r.mux.HandleFunc(path, withHeaders)
	}
	return http.ListenAndServe(port, r.mux)
}

func formatJson(data interface{}) (string, error) {
	switch v := data.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}
}

func setDefaultHeaders(rt *Router, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rt.Headers.AccessControlAllowMethods != "" {
			w.Header().Set("Access-Control-Allow-Methods", rt.Headers.AccessControlAllowMethods)
		}
		if rt.Headers.AccessControlAllowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", rt.Headers.AccessControlAllowOrigin)
		}
		if rt.Headers.ContentType != "" {
			w.Header().Set("Content-Type", rt.Headers.ContentType)
		}
		next(w, r)
	}
}

func applyMiddlewares(middlewares []func(http.HandlerFunc) http.HandlerFunc, next http.HandlerFunc) http.HandlerFunc {
	if len(middlewares) < 1 {
		return next
	}
	return middlewares[0](applyMiddlewares(middlewares[1:], next))
}
