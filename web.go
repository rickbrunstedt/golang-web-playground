package main

import (
	"app/web/router"
	"app/web/router/middlewares"
	"fmt"
	"net/http"
)

var PORT = ":8080"

type Person struct {
	Name string
	Age  int
}

func main() {
	rt := router.NewRouter()

	rt.Headers.ContentType = "text/html; charset=utf-8"
	rt.Headers.AccessControlAllowOrigin = "*"
	rt.Headers.AccessControlAllowMethods = "GET, POST, PUT, DELETE, OPTIONS"

	rt.Use(middlewares.Sessions)
	rt.Use(middlewares.Logger)

	rt.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
			<h1>Home</h1>
			<ul>
			<li><a href="/data">Data</a></li>
			<li><a href="/query?foo=bar&foo=test">Query test</a></li>
			<li><a href="/return-struct">Return struct as json</a></li>
			</ul>
		`)
	})
	rt.Get("/data", func(w http.ResponseWriter, r *http.Request) {
		data := `{"status": "ok"}`
		rt.Json(w, data)
	})
	rt.Get("/return-struct", func(w http.ResponseWriter, r *http.Request) {
		data := Person{Name: "John", Age: 30}
		rt.Json(w, data)
	})
	rt.Get("/query", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		rt.Json(w, query)
	})

	err := rt.Start(PORT)
	if err != nil {
		panic(err)
	}
}
