package main

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/Power9-Alpha/laser"
	"github.com/Power9-Alpha/laser/pkg/api"
	"github.com/Power9-Alpha/laser/pkg/postgres"

	"github.com/julienschmidt/httprouter"
)

var Start time.Time

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds | log.Llongfile)
	Start = time.Now()
}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(laser.Content, "web")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(fsys)
}

func main() {
	conf := configure()

	db, err := postgres.Init(conf.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	serviceStorage := &postgres.Service{}
	if err := serviceStorage.Init(db); err != nil {
		log.Fatal(err)
	}
	tokenStorage := &postgres.Token{}
	if err := tokenStorage.Init(db); err != nil {
		log.Fatal(err)
	}
	userStorage := &postgres.User{}
	if err := userStorage.Init(db); err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()

	router.GET("/health", Health)

	if err := api.Register(router, userStorage, tokenStorage, serviceStorage); err != nil {
		log.Fatal(err)
	}

	router.NotFound = http.FileServer(getFileSystem())

	log.Fatal(http.ListenAndServe(conf.Address, router))
}
