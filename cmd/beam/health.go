package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// @note: set from build flags via `-ldflags "-X main.Version`"
//
// @note: could add var Date string with "-X main.Date=(date +%Y%m%d-%H:%M:%S)"
var Version string = "Demo"

func Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	Status := map[string]string{"version": Version, "uptime": time.Now().Sub(Start).String()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Status)
}
