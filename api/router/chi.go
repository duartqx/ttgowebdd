package router

import (
	"html/template"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	c "github.com/duartqx/ttgowebdd/api/controllers"
	s "github.com/duartqx/ttgowebdd/application/services"
	r "github.com/duartqx/ttgowebdd/infraestructure/sqlite/repositories"
	v "github.com/duartqx/ttgowebdd/presentation/views"

	lm "github.com/duartqx/ttgowebdd/application/middleware/logger"
	rm "github.com/duartqx/ttgowebdd/application/middleware/recovery"
)

type router struct {
	db        *sqlx.DB
	templates *[]string
}

func NewRouterBuilder() *router {
	return &router{
		templates: &[]string{},
	}
}

func (ro *router) SetDb(db *sqlx.DB) *router {
	ro.db = db
	return ro
}

func (ro *router) SetTemplates(templates *[]string) *router {
	ro.templates = templates
	return ro
}

func (ro router) Build() *chi.Mux {

	tmplEngine := template.Must(template.ParseFiles(*ro.templates...))
	indexView := v.NewIndexView(tmplEngine)

	taskRepository := r.NewTaskRepository(ro.db)
	taskService := s.NewTaskService(taskRepository)

	taskController := c.NewTaskController(indexView, taskService)
	indexController := c.NewIndexController(indexView, taskService)

	router := chi.NewRouter()

	router.Use(rm.RecoveryMiddleware, lm.LoggerMiddleware)

	router.Get("/", indexController.Index)
	router.Post("/tasks/filter", indexController.Filter)

	router.Post("/tasks", taskController.Create)
	router.Put("/tasks", taskController.Update)

	return router
}
