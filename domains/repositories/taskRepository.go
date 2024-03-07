package repositories

import (
	"io"

	m "github.com/duartqx/ttgowebdd/domains/models"
)

type TaskRepository interface {
	Filter(reader io.ReadCloser) (*[]m.Task, error)

	FindById(id int64) (m.Task, error)
	FindByTag(tag string) (m.Task, error)

	Create(task m.Task) error

	Complete(task m.Task) error
	CompleteById(id int) (m.Task, error)

	UpdateEndAt(task m.Task) error
	UpdateEndAtById(id int) (m.Task, error)

	GetListOfTaskSprints() *[]string
}
