package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "ini /")
	})

	http.HandleFunc("/index", index)

	fmt.Println("starting web server at http://localhost:8080")

	/*
		membuka web server pada port 8080
		fungsi ini menggunakan goroutine
		agar dapat berjalan di background
	*/
	http.ListenAndServe(":8080", nil)
}
