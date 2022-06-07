package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Auth struct {
	user  authStorage
	token tokenStorage
}

func (a *Auth) hasBearer(r *http.Request) (string, bool) {
	authorization := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", false
	}
	return strings.TrimPrefix(authorization, "Bearer "), true
}

// @note: ideally we would have scope as a GET parameter to either of these
//        operations, which would be used to add permissions to the JWT claims
//        provided in a successful response.

// @note: currently we have nothing that automates creation of refresh tokens
//        so we may need to add some behavior to the project that either
//        automates token creation, or expects a request to create a token
//        without authentication, as the response does not include the token
//        and the identified is unique, making it "safe" to call without
//        permissions.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var bearer string
	var hasBearer bool
	if username, password, hasBasicAuth := r.BasicAuth(); hasBasicAuth {
		if user, err := a.user.Login(username, password); err != nil {
			log.Println("failed to authenticate (%s): %s", username, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if token, err := a.token.Select(user.Name); err != nil {
			log.Println("failed to acquire token by username (%s): %s", user.Name, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else {
			bearer = token.ID
		}
	} else if bearer, hasBearer = a.hasBearer(r); !hasBearer {
		log.Println("Missing Authorization...")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// @note: there may be a cleaner way to write the above block, a problem
	//        for future iterations.

	reply := map[string]any{}
	reply["token_type"] = "Bearer"
	reply["refresh_token"] = bearer
	// @todo: generate and return JWT /w refresh token as json response
	reply["access_token"] = ""

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(reply)
}

func (a *Auth) Access(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bearer, hasBearer := a.hasBearer(r)
	if !hasBearer {
		log.Println("Missing Bearer Authorization...")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// @todo: add redis cache bearer check to rate-limit access token requests

	if _, err := a.token.Select(bearer); err != nil {
		log.Println("Failed to acquire bearer　token (%s): %s", bearer, err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	reply := map[string]any{}
	reply["token_type"] = "Bearer"
	reply["refresh_token"] = bearer
	// @todo: generate and return JWT /w refresh token as json response
	reply["access_token"] = ""

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(reply)
}

// @todo: add wrapper function to check jwt for matching permissions, and
//        which also leverages cors wrapper
