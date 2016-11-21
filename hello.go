package main

import (
	"fmt"
	"net/http"

	"github.com/larshelmer/hello/storage"
	"github.com/larshelmer/hello/v1_api"
)

// auth0
// kibana
// https://elithrar.github.io/article/testing-http-handlers-go/

func oldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello %s", r.URL.Path[1:])
}

func motdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "gurka är gott")
}

func main() {
	fmt.Println("starting...")

	d := storage.Storage{}

	d.InitData("")
	v1api.InitEndpoints(d)

	http.Handle("/", http.FileServer(http.Dir("static")))
	//	http.HandleFunc("/old", oldHandler)
	//	http.HandleFunc("/motd", motdHandler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("stopping...")
}
