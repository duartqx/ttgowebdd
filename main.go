package main

import (
	"net/http"
	"time"

	"github.com/duartqx/ttgowebdd/api/router"
	"github.com/duartqx/ttgowebdd/infraestructure/sqlite"
)

func main() {
	db, err := sqlite.NewDbConnection("ddclitskt.sqlite")
	if err != nil {
		panic(err)
	}
	mux := router.NewRouterBuilder().
		SetDb(db).
		SetTemplates(&[]string{
			"./presentation/templates/createForm.html",
			"./presentation/templates/filterForm.html",
			"./presentation/templates/formTabs.html",
			"./presentation/templates/index.html",
			"./presentation/templates/taskTable.html",
		}).
		Build()

	port := ":8000"

	srv := &http.Server{
		Handler:      mux,
		Addr:         "127.0.0.1" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
