package main

import (
	"app/web/router"
	"app/web/router/middlewares"
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
		<main>
			<h1>Home</h1>
			<ul>
				<li><a href="/data">Data</a></li>
				<li><a href="/query?foo=bar&foo=test">Query test</a></li>
				<li><a href="/return-struct">Return struct as json</a></li>
				<li><a href="/test-post">Test POST</a></li>
				<li><a href="/dynamic-params/123">Dynamic params</a></li>
				<li><a href="/dynamic-params/123/nested/bar">Dynamic params</a></li>
			</ul>

			<div>
				<h2>Testing variables in template</h2>
				<p>Name: {{.Name}}</p>
				<p>Age: {{.Age}}</p>
			</div>
		</main>
		`
		ctx.Html(html, Person{Name: "Bob", Age: 67})
	})
	rt.Get("/dynamic-params/:foo", func(ctx router.RouteCtx) {
		someId := ctx.Params["foo"]
		html := `
			<main>
				<h1>Dynamic params</h1>
				<p>someId: {{.}}</p>
			</main>
		`
		ctx.Html(html, someId)
	})
	rt.Get("/dynamic-params/:foo/nested/:bar", func(ctx router.RouteCtx) {
		html := `
			<main>
				<h1>Dynamic params</h1>
				<p>Foo: {{.foo}}</p>
				<p>Bar: {{.bar}}</p>
			</main>
		`
		ctx.Html(html, ctx.Params)
	})
	rt.Get("/data", func(ctx router.RouteCtx) {
		data := `{"status": "ok"}`
		ctx.Json(data)
	})
	rt.Get("/return-struct", func(ctx router.RouteCtx) {
		data := Person{Name: "Bob", Age: 67}
		ctx.Json(data)
	})
	rt.Get("/query", func(ctx router.RouteCtx) {
		query := ctx.R.URL.Query()
		ctx.Json(query)
	})
	rt.Get("/test-post", func(ctx router.RouteCtx) {
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
		name := ctx.R.FormValue("name")
		html := `
			<main>
				<h1>Test POST</h1>
				<p>Hello {{.}}</p>
			</main>
		`
		ctx.Html(html, name)
	})

	log.Println("Listening on port", PORT)

	err := rt.Start(PORT)
	if err != nil {
		panic(err)
	}
}
