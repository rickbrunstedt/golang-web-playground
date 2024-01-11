package main

import (
	"app/web/router"
	"app/web/router/middlewares"
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

	// rt.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprint(w, `
	// 		<h1>Home</h1>
	// 		<ul>
	// 		<li><a href="/data">Data</a></li>
	// 		<li><a href="/query?foo=bar&foo=test">Query test</a></li>
	// 		<li><a href="/return-struct">Return struct as json</a></li>
	// 		<li><a href="/test-post">Test POST</a></li>
	// 		</ul>
	// 	`)
	// })
	// rt.Get("/data", func(w http.ResponseWriter, r *http.Request) {
	// 	data := `{"status": "ok"}`
	// 	rt.Json(w, data)
	// })
	// rt.Get("/return-struct", func(w http.ResponseWriter, r *http.Request) {
	// 	data := Person{Name: "John", Age: 30}
	// 	rt.Json(w, data)
	// })
	// rt.Get("/query", func(w http.ResponseWriter, r *http.Request) {
	// 	query := r.URL.Query()
	// 	rt.Json(w, query)
	// })
	// rt.Get("/test-post", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("GET /test-post")
	// 	fmt.Fprint(w, `
	// 		<main>
	// 			<h1>Test POST</h1>
	// 			<form action="/test-post" method="post">
	// 				<input type="text" name="name" />
	// 				<input type="submit" value="Submit" />
	// 			</form>
	// 		</main>
	// 	`)
	// })
	// rt.Post("/test-post", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("POST /test-post")
	// 	name := r.FormValue("name")
	// 	fmt.Fprintf(w, "Hello %s", name)
	// })

	err := rt.Start(PORT)
	if err != nil {
		panic(err)
	}
}
