package services

import (
	"encoding/json"
	"io"

	t "github.com/duartqx/ttgowebdd/domains/entities/task"
	m "github.com/duartqx/ttgowebdd/domains/models"
	r "github.com/duartqx/ttgowebdd/domains/repositories"
)

type TaskService struct {
	repo r.TaskRepository
}

func NewTaskService(repo r.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (ts TaskService) GetTasksByFilter(reader io.ReadCloser) (*[]m.Task, error) {
	return ts.repo.Filter(reader)
}

func (ts TaskService) GetListOfTaskSprints() *[]string {
	return ts.repo.GetListOfTaskSprints()
}

func (ts TaskService) Create(reader io.ReadCloser) (m.Task, error) {
	var task t.TaskEntity
	if err := json.NewDecoder(reader).Decode(&task); err != nil {
		return nil, err
	}
	return &task, ts.repo.Create(&task)
}

func (ts TaskService) UpdateEndAt(id int) (m.Task, error) {
	return ts.repo.UpdateEndAtById(id)
}

func (ts TaskService) UpdateComplete(id int) (m.Task, error) {
	return ts.repo.CompleteById(id)
}

func (ts TaskService) GetAll() (*[]m.Task, error) {
	return ts.repo.Filter(nil)
}
