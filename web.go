package main

import (
	"app/web/router"
	"encoding/json"
	"fmt"
	"net/http"
)

var PORT = ":8080"

type Person struct {
	name string
	age  int
}

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
			<li><a href="/return-struct">Return struct as json</a></li>
			</ul>
		`)
	})
	rt.Get("/data", func(w http.ResponseWriter, r *http.Request) {
		data := `{"status": "ok"}`
		rt.Json(w, data)
	})
	rt.Get("/return-struct", func(w http.ResponseWriter, r *http.Request) {
		data := Person{name: "John", age: 30}
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		rt.Json(w, jsonData)
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
