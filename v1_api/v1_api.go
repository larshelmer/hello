package v1api

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/larshelmer/hello/storage"
)

type env struct {
	db storage.Datastore
}

var (
	messageURL       = "/v1/message/"
	randomMessageURL = "/v1/message/random"
)

// InitEndpoints initializes handlers for all endpoints
func InitEndpoints(s storage.Datastore) {
	env := &env{s}
	http.HandleFunc(randomMessageURL, env.getRandomMessageHandler)
	http.HandleFunc(messageURL, env.messageHandler)
}

func (e *env) getAllMessages(w http.ResponseWriter) {
	dat, err := e.db.Read()
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
}

func (e *env) getOneMessage(w http.ResponseWriter, r *http.Request, id string) {
	i, err := strconv.Atoi(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if i < 0 {
		http.NotFound(w, r)
		return
	}
	dat, err := e.db.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if i < len(*dat) {
		js, err := json.Marshal((*dat)[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	} else {
		http.NotFound(w, r)
	}
}

func (e *env) messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Path[len(messageURL):]
		if len(id) == 0 {
			e.getAllMessages(w)
		} else {
			e.getOneMessage(w, r, id)
		}
	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var motd string
		err := decoder.Decode(&motd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer r.Body.Close()
		err = e.db.Add(motd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (e *env) getRandomMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	dat, err := e.db.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if len(*dat) == 0 {
			w.WriteHeader(http.StatusNoContent)
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
}
