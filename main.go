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

type config struct {
	dbStr              string
	hostAddr           string
	jiraAuth           string
	jiraSearchEndpoint string
	jiraAssignee       string
}

func GetConfig() *config {
	return &config{
		dbStr:              os.Getenv("dbStr"),
		hostAddr:           os.Getenv("hostAddr"),
		jiraAuth:           os.Getenv("jiraAuth"),
		jiraSearchEndpoint: os.Getenv("jiraSearchEndpoint"),
		jiraAssignee:       os.Getenv("jiraAssignee"),
	}
}
