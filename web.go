package main

import (
	"app/web/router"
	"fmt"
	"net/http"
)

var PORT = ":8080"

func main() {
	rt := router.NewRouter()

	rt.Use(router.HandleSession)
	rt.Use(router.LoggerMiddleware)

	rt.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<h1>Home</h1>
			<ul>
			<li><a href="/data">Data</a></li>
			<li><a href="/query?foo=bar&foo=test">Query test</a></li>
			</ul>
		`)
	})
	rt.Get("/data", func(w http.ResponseWriter, r *http.Request) {
		data := `{"status": "ok"}`
		rt.Json(w, data)
	})
	rt.Get("/query", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		rt.Json(w, query)
	})

	rt.Start(PORT)
}
