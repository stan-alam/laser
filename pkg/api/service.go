package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Power9-Alpha/laser"

	"github.com/julienschmidt/httprouter"
)

type Service struct {
	storage serviceStorage
}

func (h *Service) Post(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var service laser.Service

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&service); err != nil {
		log.Printf("failed to decode json: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if h.storage.Insert(&service) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Service) GetOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	service, err := h.storage.SelectOne(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if service == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(service)
}

func (h *Service) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	services, err := h.storage.Select()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(services)
}

func (h *Service) Put(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var service laser.Service

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&service); err != nil {
		log.Printf("failed to decode json: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.storage.Update(&service); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Service) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if h.storage.Delete(id) != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
