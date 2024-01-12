package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []func(http.HandlerFunc) http.HandlerFunc
	Headers     Headers
	routes      map[string]Route
}

type Route struct {
	Path       string
	PathReg    *regexp.Regexp
	ParamNames []string
	Handler    func(http.ResponseWriter, *http.Request, map[string]string)
	Method     string
}

// Maybe this shouldn't be name Ctx because it's not a context as in context.Context
type RouteCtx struct {
	W        http.ResponseWriter
	R        *http.Request
	Redirect func(string, int)
	Json     func(interface{})
	Html     func(string, ...interface{})
	Params   map[string]string
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

func (r *Router) addRoute(path string, method string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	pathRegex, dynamicNames := makePathRegexp(path)
	routeKey := fmt.Sprintf("%s:%s", path, method)
	r.routes[routeKey] = Route{
		Path:       path,
		PathReg:    pathRegex,
		ParamNames: dynamicNames,
		Handler:    handler,
		Method:     method,
	}
}

func (rt *Router) makeHandler(handler func(RouteCtx)) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		handler(RouteCtx{
			W: w,
			R: r,
			Redirect: func(path string, status int) {
				http.Redirect(w, r, path, status)
			},
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
			Params: params,
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
	path := r.URL.Path

	for _, route := range rt.routes {
		if route.PathReg.MatchString(path) {
			params := extractParamsFromURL(route, path)
			route.Handler(w, r, params)
			return
		}
	}

	http.Error(w, "Not found", 404)
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

func makePathRegexp(path string) (*regexp.Regexp, []string) {
	regexRoute := regexp.MustCompile(`:[^/]+`).ReplaceAllString(path, `[^/]+`)
	regexMatches := regexp.MustCompile(`:([^/]+)`).FindAllStringSubmatch(path, -1)
	matches := make([]string, len(regexMatches))
	for i, match := range regexMatches {
		matches[i] = match[1]
	}
	regexRoute = "^" + regexRoute + "$"
	fullRegex := regexp.MustCompile(regexRoute)
	return fullRegex, matches
}

func extractParamsFromURL(route Route, path string) map[string]string {
	params := make(map[string]string)

	// Split the route's path and the incoming path into segments
	routeSegments := strings.Split(route.Path, "/")
	pathSegments := strings.Split(path, "/")

	// Ensure the number of segments match
	if len(routeSegments) == len(pathSegments) {
		for i, segment := range routeSegments {
			if strings.HasPrefix(segment, ":") && i < len(pathSegments) {
				paramName := strings.TrimPrefix(segment, ":")
				params[paramName] = pathSegments[i]
			}
		}
	}

	return params
}
