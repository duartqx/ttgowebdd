package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	s "github.com/duartqx/ttgowebdd/application/services"
	v "github.com/duartqx/ttgowebdd/presentation/views"
)

type IndexController struct {
	view    *v.IndexView
	service *s.TaskService
}

func NewIndexController(view *v.IndexView, service *s.TaskService) *IndexController {
	return &IndexController{
		view:    view,
		service: service,
	}
}

func (ic IndexController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ic.Index(w, r)
	case http.MethodPost:
		ic.Filter(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (ic IndexController) Index(w http.ResponseWriter, r *http.Request) {
	sprints := ic.service.GetListOfTaskSprints()
	data := map[string]interface{}{
		"Sprints": sprints,
	}
	log.Println(sprints)
	if err := ic.view.Execute(w, data); err != nil {
		panic(err)
	}
}

func (ic IndexController) Filter(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	formData, err := json.Marshal(r.Form)
	if err != nil {
		panic(err)
	}

	tasks, err := ic.service.GetTasksByFilter(io.NopCloser(bytes.NewReader(formData)))
	if err != nil {
		if strings.HasPrefix(err.Error(), "Decode Error:") {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			panic(err)
		}
	}
	data := map[string]interface{}{
		"Tasks": tasks,
	}
	if err := ic.view.ExecuteResults(w, data); err != nil {
		panic(err)
	}
}
