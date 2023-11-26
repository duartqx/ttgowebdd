package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	s "github.com/duartqx/ttgowebdd/application/services"
	v "github.com/duartqx/ttgowebdd/presentation/views"
)

type TaskController struct {
	view    *v.IndexView
	service *s.TaskService
}

func NewTaskController(view *v.IndexView, service *s.TaskService) *TaskController {
	return &TaskController{
		view:    view,
		service: service,
	}
}

func (tc TaskController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		tc.Create(w, r)
	case http.MethodPut:
		tc.UpdateEndAt(w, r)
	case http.MethodPatch:
		tc.UpdateComplete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (tc TaskController) getId(r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (tc TaskController) UpdateEndAt(w http.ResponseWriter, r *http.Request) {

	id, err := tc.getId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := tc.service.UpdateEndAt(id)
	if err != nil {
		panic(err)
	}
	if err := tc.view.ExecuteRow(w, task); err != nil {
		panic(err)
	}
}

func (tc TaskController) UpdateComplete(w http.ResponseWriter, r *http.Request) {
	id, err := tc.getId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := tc.service.UpdateComplete(id)
	if err != nil {
		panic(err)
	}
	if err := tc.view.ExecuteRow(w, task); err != nil {
		panic(err)
	}
}

func (tc TaskController) Create(w http.ResponseWriter, r *http.Request) {
	task, err := tc.service.Create(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := tc.view.ExecuteRow(w, task); err != nil {
		panic(err)
	}
}
