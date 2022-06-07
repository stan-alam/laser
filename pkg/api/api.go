package api

import (
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
func Register(router *httprouter.Router, a authStorage, t tokenStorage, s serviceStorage) error {
	router.GlobalOPTIONS = http.HandlerFunc(CORSOptions)

	// @note: permission wrapper using jwt
	// should be able to check permissions there
	// such as `admin` scope, exceptions being data
	// ownership, such as a user deleting themselves
	// which we may not want to allow anyways...
	// router.POST("/oauth/token", userHandler.Login) // @note: does this belong attached to user?

	// // @todo: once auth is updated we can apply permission check wrappers
	// router.GET("/api/user/:id", CorsWrapper(userHandler.GetOne))
	// router.DELETE("/api/user/:id", CorsWrapper(userHandler.Delete))
	// router.POST("/api/user", CorsWrapper(userHandler.Post))
	// router.GET("/api/users", CorsWrapper(userHandler.Get))
	// router.PUT("/api/user/:id", CorsWrapper(userHandler.Put))

	// router.GET("/api/service/:id", CorsWrapper(serviceHandler.GetOne))
	// router.DELETE("/api/service/:id", CorsWrapper(serviceHandler.Delete))
	// router.POST("/api/service", CorsWrapper(serviceHandler.Post))
	// router.PUT("/api/service", CorsWrapper(serviceHandler.Put))
	// router.GET("/api/services", CorsWrapper(serviceHandler.Get))

	return nil
}
