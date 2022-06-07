package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Power9-Alpha/laser"

	"github.com/julienschmidt/httprouter"
)

type Token struct {
	Storage tokenStorage
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

	if _, err := t.Storage.Insert(token.Username); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *Token) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	token, err := t.Storage.Select(id)
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
	if t.Storage.Delete(id) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
