package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []func(http.HandlerFunc) http.HandlerFunc
	Headers     Headers
	routes      map[string]Route
}

type Route struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
	Method  string
}

// Maybe this shouldn't be name Ctx because it's not a context as in context.Context
type RouteCtx struct {
	W    http.ResponseWriter
	R    *http.Request
	Json func(interface{})
	Html func(string, ...interface{})
}

type RouteHandler func(RouteCtx)

type Headers struct {
	ContentType               string
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
}

func NewRouter() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		middlewares: []func(http.HandlerFunc) http.HandlerFunc{},
		Headers:     Headers{},
		routes:      map[string]Route{},
	}
}

func (r *Router) addRoute(path string, method string, handler func(http.ResponseWriter, *http.Request)) {
	routeKey := fmt.Sprintf("%s:%s", path, method)
	r.routes[routeKey] = Route{
		Path:    path,
		Handler: handler,
		Method:  method,
	}
}

func (rt *Router) makeHandler(handler func(RouteCtx)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(RouteCtx{
			W: w,
			R: r,
			Json: func(i interface{}) {
				rt.Json(w, i)
			},
			Html: func(html string, data ...interface{}) {
				if len(data) < 1 {
					rt.Html(w, html, nil)
				} else {
					rt.Html(w, html, data[0])
				}
			},
		})
	}
}

func (r *Router) Get(path string, handler func(RouteCtx)) {
	wrappedHandler := r.makeHandler(handler)
	r.addRoute(path, http.MethodGet, wrappedHandler)
}

func (r *Router) Post(path string, handler func(RouteCtx)) {
	wrappedHandler := r.makeHandler(handler)
	r.addRoute(path, http.MethodPost, wrappedHandler)
}

func (r *Router) Put(path string, handler func(RouteCtx)) {
	wrappedHandler := r.makeHandler(handler)
	r.addRoute(path, http.MethodPut, wrappedHandler)
}

func (r *Router) Delete(path string, handler func(RouteCtx)) {
	wrappedHandler := r.makeHandler(handler)
	r.addRoute(path, http.MethodDelete, wrappedHandler)
}

func (r *Router) Options(path string, handler func(RouteCtx)) {
	wrappedHandler := r.makeHandler(handler)
	r.addRoute(path, http.MethodOptions, wrappedHandler)
}

// How would this work? Would it be a middleware?
func (r *Router) Head(path string, handler func(RouteCtx)) {
	wrappedHandler := r.makeHandler(handler)
	r.addRoute(path, http.MethodHead, wrappedHandler)
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

func (r *Router) Html(w http.ResponseWriter, html string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.New("random").Parse(html)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func (rt *Router) routesHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	routeKey := fmt.Sprintf("%s:%s", path, method)
	route, ok := rt.routes[routeKey]
	if !ok {
		http.Error(w, "Not found", 404)
		return
	}
	route.Handler(w, r)
}

func (r *Router) Start(port string) error {
	routingWithMiddlewares := applyMiddlewares(r.middlewares, r.routesHandler)
	routingWithHeaders := setDefaultHeaders(r, routingWithMiddlewares)
	r.mux.HandleFunc("/", routingWithHeaders)
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

// TODO: This should be done in a better way
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
