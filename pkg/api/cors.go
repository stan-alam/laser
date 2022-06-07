package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// @note: this seems to be required by the browser on non-OPTION requests
func CorsWrapper(h httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		h(w, r, ps)
	})
}

func CORSOptions(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Access-Control-Request-Method") != "" {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Allow", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Allow")) // alternative: "GET,PUT,POST,PATCH,DELETE,OPTIONS"
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.Header().Set("Access-Control-Allow-Origin", origin) // @note: "*" sometimes does not work
		w.Header().Set("Vary", "Origin")
	}

	// @note: alternatively could send http.StatusNoContent (204), but then must omit headers
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
