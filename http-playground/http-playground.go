package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := "test 1"
		fmt.Println(data)
		fmt.Fprint(w, data)
	})
	// mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
	// 	data := "test 2"
	// 	fmt.Println(data)
	// 	fmt.Fprint(w, data)
	// })
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// data := `{"status": "ok"}`
	// 	data := "test 2"
	// 	fmt.Println(data)
	// 	fmt.Fprint(w, data)
	// })

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
