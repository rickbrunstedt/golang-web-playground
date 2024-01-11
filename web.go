package main

import (
	"app/web/router"
	"app/web/router/middlewares"
	"fmt"
	"log"
)

var PORT = ":8080"

type Person struct {
	Name string
	Age  int
}

func main() {
	rt := router.NewRouter()

	// Setting default headers
	rt.Headers.ContentType = "text/html; charset=utf-8"
	rt.Headers.AccessControlAllowOrigin = "*"
	rt.Headers.AccessControlAllowMethods = "GET, POST, PUT, DELETE, OPTIONS"

	rt.Use(middlewares.Sessions)
	rt.Use(middlewares.Logger)

	rt.Get("/", func(ctx router.RouteCtx) {
		html := `
			<h1>Home</h1>
			<ul>
			<li><a href="/data">Data</a></li>
			<li><a href="/query?foo=bar&foo=test">Query test</a></li>
			<li><a href="/return-struct">Return struct as json</a></li>
			<li><a href="/test-post">Test POST</a></li>
			</ul>
		`
		ctx.Html(html)
	})

	rt.Get("/data", func(ctx router.RouteCtx) {
		data := `{"status": "ok"}`
		ctx.Json(data)
	})
	rt.Get("/return-struct", func(ctx router.RouteCtx) {
		data := Person{Name: "John", Age: 30}
		ctx.Json(data)
	})
	rt.Get("/query", func(ctx router.RouteCtx) {
		query := ctx.R.URL.Query()
		ctx.Json(query)
	})
	rt.Get("/test-post", func(ctx router.RouteCtx) {
		log.Println("GET /test-post")
		html := `
			<main>
				<h1>Test POST</h1>
				<form action="/test-post" method="post">
					<input type="text" name="name" />
					<input type="submit" value="Submit" />
				</form>
			</main>
		`
		ctx.Html(html)
	})
	rt.Post("/test-post", func(ctx router.RouteCtx) {
		log.Println("POST /test-post")
		name := ctx.R.FormValue("name")
		ctx.W.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(ctx.W, "Hello %s", name)
	})

	log.Println("Listening on port", PORT)

	err := rt.Start(PORT)
	if err != nil {
		panic(err)
	}
}
