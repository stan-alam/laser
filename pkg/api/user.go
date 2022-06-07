package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Power9-Alpha/laser"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	storage userStorage
}

func (h *User) Post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user laser.User

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&user); err != nil {
		log.Printf("failed to decode json: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// @note: alternative would be direct form parsing
	// if err := r.ParseForm(); err != nil {
	// 	log.Printf("Failed to parse form submission: %s", err)
	// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }

	// user := laser.User{
	// 	Name:     r.Form.Get("username"),
	// 	Email:    r.Form.Get("email"),
	// 	Password: r.Form.Get("password"),
	// }

	if h.storage.Insert(&user) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *User) GetOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	user, err := h.storage.SelectOne(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if user == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(user)

}

func (h *User) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users, err := h.storage.Select()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(users)
}

func (h *User) Put(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user laser.User

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&user); err != nil {
		log.Printf("failed to decode json: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.storage.Update(&user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *User) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if h.storage.Delete(id) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
