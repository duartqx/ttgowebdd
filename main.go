package main

import (
	"embed"
	"net/http"
	"os"
	"time"

	"github.com/duartqx/ttgowebdd/api/router"
	"github.com/duartqx/ttgowebdd/infraestructure/sqlite"
)

//go:embed presentation/templates/*
var tmplFolder embed.FS

func main() {
	config := GetConfig()

	db, err := sqlite.NewDbConnection(config.dbStr)
	if err != nil {
		panic(err)
	}
	mux := router.NewRouterBuilder().
		SetDb(db).
		SetTemplFolder(&tmplFolder).
		SetTemplates(&[]string{
			"presentation/templates/*.html",
		}).
		Build()

	srv := &http.Server{
		Handler:      mux,
		Addr:         config.hostAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

type Config struct {
	dbStr    string
	hostAddr string
}

func GetConfig() *Config {
	config := Config{
		dbStr:    os.Getenv("dbStr"),
		hostAddr: os.Getenv("hostAddr"),
	}

	if config.dbStr == "" {
		panic("Missing dbStr")
	}
	if config.hostAddr == "" {
		panic("Missing hostAddr")
	}

	return &config
}
