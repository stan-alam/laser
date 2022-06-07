package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Power9-Alpha/laser"

	"github.com/julienschmidt/httprouter"
)

type Token struct {
	storage tokenStorage
}

func (t *Token) Post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := &laser.Token{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(token); err != nil {
		log.Printf("failed to decode json: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if _, err := t.storage.Insert(token.Username); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *Token) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	token, err := t.storage.Select(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(token)
}

func (t *Token) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if t.storage.Delete(id) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
