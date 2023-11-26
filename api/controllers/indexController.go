package controllers

import (
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
	tasks, err := ic.service.GetAll()
	if err != nil {
		// Recover middleware will catch this panic
		panic(err)
	}
	if err := ic.view.Execute(w, tasks); err != nil {
		panic(err)
	}
}

// Filter is a method of the IndexController struct used to filter tasks based on various criteria.
// This method accepts an HTTP request and an HTTP response writer as arguments.
// It reads the request body and passes it to the GetTasksByFilter service to retrieve tasks matching the filter criteria.
// The filter criteria are supplied in the request body in JSON format. The following fields are supported:
// - "tag": A string specifying the tag of the tasks to retrieve.
// - "completed": An integer specifying the completion status of the tasks to retrieve. Can be one of the following:
//   - 0: Ignore the completion status.
//   - 1: Retrieve only completed tasks.
//   - 2: Retrieve only not completed tasks.
//
// - "from_start_at": A time.Time specifying the earliest start time of the tasks to retrieve.
// - "to_start_at": A time.Time specifying the latest start time of the tasks to retrieve.
// - "from_end_at": A time.Time specifying the earliest end time of the tasks to retrieve.
// - "to_end_at": A time.Time specifying the latest end time of the tasks to retrieve.
// If an error occurs while decoding the request body, the method responds with a 400 Bad Request status.
// If any other error occurs while retrieving the tasks, the method panics.
// If the tasks are retrieved successfully, they are passed to the Execute method of the view along with the response writer to generate the response.
// If an error occurs while executing the view, the method panics.
func (ic IndexController) Filter(w http.ResponseWriter, r *http.Request) {
	tasks, err := ic.service.GetTasksByFilter(r.Body)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Decode Error:") {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			panic(err)
		}
	}
	if err := ic.view.ExecuteResults(w, tasks); err != nil {
		panic(err)
	}
}
