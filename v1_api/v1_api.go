package v1api

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/larshelmer/hello/storage"
)

// InitEndpoints initializes handlers for all endpoints
func InitEndpoints() {
	http.HandleFunc("/v1/message/random", getRandomMessageHandler)
	http.HandleFunc("/v1/message", messageHandler)
}

func getOpenAPIDefinition(w http.ResponseWriter, r *http.Request) {

}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		dat, err := storage.Read()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			js, err := json.Marshal(*dat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
		}
	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var motd string
		err := decoder.Decode(&motd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer r.Body.Close()
		err = storage.Add(motd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func getRandomMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	dat, err := storage.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		ix := rand.Int() % len(*dat)
		js, err := json.Marshal((*dat)[ix])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}
