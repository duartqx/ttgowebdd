package controllers

import (
	"encoding/json"
	"net/http"

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
		tc.Update(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (tc TaskController) Update(w http.ResponseWriter, r *http.Request) {
	const (
		updEnd      string = "end"
		updComplete        = "complete"
	)

	putBody := struct {
		Id        int    `json:"id"`
		Operation string `json:"operation"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&putBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		task interface{}
		err  error
	)

	switch putBody.Operation {
	case updEnd:
		task, err = tc.service.UpdateEndAt(putBody.Id)
	case updComplete:
		task, err = tc.service.UpdateComplete(putBody.Id)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

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
