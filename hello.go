package main

import (
	"fmt"
	"net/http"

	"github.com/larshelmer/hello/storage"
)

func oldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %s", r.URL.Path[1:])
	storage.Read()
}

func motdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "gurka är gott")
}

func main() {
	fmt.Println("starting...")
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/old", oldHandler)
	http.HandleFunc("/motd", motdHandler)
	http.ListenAndServe(":8080", nil)
}
