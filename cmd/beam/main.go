package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

var Version string = "Demo"
var Start time.Time

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds | log.Llongfile)
}

func Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	Status := map[string]string{"version": Version, "uptime": time.Now().Sub(Start).String()}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Status)
}

func main() {
	Start = time.Now()

	conf := configure()

	router := httprouter.New()
	router.GET("/health", Health)

	log.Fatal(http.ListenAndServe(conf.Address, router))
}
