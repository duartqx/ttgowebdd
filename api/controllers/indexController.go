package controllers

import (
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
	tasks, err := ic.service.GetAll()
	if err != nil {
		// Recover middleware will catch this panic
		panic(err)
	}
	if err := ic.view.Execute(w, tasks); err != nil {
		panic(err)
	}
}

func (ic IndexController) Filter(w http.ResponseWriter, r *http.Request) {
	tasks, err := ic.service.GetTasksByFilter(r.Body)
	if err != nil {
		if strings.HasPrefix("Decode Error:", err.Error()) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			panic(err)
		}
	}
	if err := ic.view.Execute(w, tasks); err != nil {
		panic(err)
	}
}
