package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func main() {
	// nampilin ini di localhost:8080
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ini /")
	})

	// kalo localhost:8080/index nnt nampilin yang ada di fungsi index
	// handleFunc, nambahin fungsi di parameter pertama
	http.HandleFunc("/index", index)

	// nunjukin kalo jalan
	fmt.Println("starting web server at http://localhost:8080/")

	// listen dan serve di port 8080, ini pake goroutine
	http.ListenAndServe(":8080", nil)
}
