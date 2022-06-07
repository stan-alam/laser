package api

import (
	"net/http"

	"github.com/Power9-Alpha/laser"

	"github.com/julienschmidt/httprouter"
)

type userStorage interface {
	Insert(*laser.User) error
	SelectOne(string) (*laser.User, error)
	Select() ([]laser.User, error)
	Update(*laser.User) error
	Delete(string) error
}

type serviceStorage interface {
	Insert(*laser.Service) error
	SelectOne(string) (*laser.Service, error)
	Select() ([]laser.Service, error)
	Update(*laser.Service) error
	Delete(string) error
}

type tokenStorage interface {
	Insert(string) (string, error)
	Select(string) (*laser.Token, error)
	Delete(string) error
}

type authStorage interface {
	Login(string, string) (*laser.User, error)
}

// @note: composite interface reduces Register method signature
type authUserStorage interface {
	authStorage
	userStorage
}

// @note: Register method accepts an httprouter, and all interfaces,
//        and applies inversion-of-control to initialize and register routes
//        for all operations, returning an error.
func Register(router *httprouter.Router, a authUserStorage, t tokenStorage, s serviceStorage) error {
	users := User{storage: a}
	tokens := Token{storage: t}
	services := Service{storage: s}
	auth := Auth{user: a, token: t}

	router.GlobalOPTIONS = http.HandlerFunc(CORSOptions)

	// @note: we may want to prefix these with `/api`?
	router.POST("/oauth/token", auth.Login)
	router.POST("/oauth/Access", auth.Access)

	router.GET("/api/user/:id", CorsWrapper(users.GetOne))
	router.DELETE("/api/user/:id", CorsWrapper(users.Delete))
	router.POST("/api/user", CorsWrapper(users.Post))
	router.GET("/api/users", CorsWrapper(users.Get))
	router.PUT("/api/user/:id", CorsWrapper(users.Put))

	router.GET("/api/service/:id", CorsWrapper(services.GetOne))
	router.DELETE("/api/service/:id", CorsWrapper(services.Delete))
	router.POST("/api/service", CorsWrapper(services.Post))
	router.PUT("/api/service", CorsWrapper(services.Put))
	router.GET("/api/services", CorsWrapper(services.Get))

	router.GET("/api/token/:id", CorsWrapper(tokens.Get))
	router.DELETE("/api/token/:id", CorsWrapper(tokens.Delete))
	router.POST("/api/token", CorsWrapper(tokens.Post))

	return nil
}
