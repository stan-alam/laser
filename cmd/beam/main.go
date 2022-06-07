package main

import (
	"database/sql"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/Power9-Alpha/laser"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
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

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(laser.Content, "web")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(fsys)
}

func main() {
	Start = time.Now()

	conf := configure()

	db, err := sql.Open("postgres", conf.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(40)

	userStorage := UserStorage{db: db}
	serviceStorage := ServiceStorage{db: db}
	if err := userStorage.Init(); err != nil {
		log.Fatal(err)
	} else if err = serviceStorage.Init(); err != nil {
		log.Fatal(err)
	}

	userHandler := &UserHandler{storage: userStorage}
	serviceHandler := &ServiceHandler{storage: serviceStorage}

	router := httprouter.New()

	// register RESTful operations /w Options for CORS in development/local
	// router.OPTIONS("/*all", CORSOptions) // newer httprouter recommended GlobalOPTIONS
	router.GlobalOPTIONS = http.HandlerFunc(CORSOptions)

	router.GET("/health", Health)

	// @note: permission wrapper using jwt
	// should be able to check permissions there
	// such as `admin` scope, exceptions being data
	// ownership, such as a user deleting themselves
	// which we may not want to allow anyways...
	// router.POST("/oauth/token", userHandler.Login) // @note: does this belong attached to user?

	// @todo: once auth is updated we can apply permission check wrappers
	router.GET("/api/user/:id", CorsWrapper(userHandler.GetOne))
	router.DELETE("/api/user/:id", CorsWrapper(userHandler.Delete))
	router.POST("/api/user", CorsWrapper(userHandler.Post))
	router.GET("/api/users", CorsWrapper(userHandler.Get))
	router.PUT("/api/user/:id", CorsWrapper(userHandler.Put))

	router.GET("/api/service/:id", CorsWrapper(serviceHandler.GetOne))
	router.DELETE("/api/service/:id", CorsWrapper(serviceHandler.Delete))
	router.POST("/api/service", CorsWrapper(serviceHandler.Post))
	router.PUT("/api/service", CorsWrapper(serviceHandler.Put))
	router.GET("/api/services", CorsWrapper(serviceHandler.Get))

	router.NotFound = http.FileServer(getFileSystem())

	log.Fatal(http.ListenAndServe(conf.Address, router))
}
