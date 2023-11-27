package task

import (
	"time"

	m "github.com/duartqx/ttgowebdd/domains/models"
)

type TaskEntity struct {
	Id          int64      `json:"id" db:"id"`
	Tag         string     `json:"tag" db:"tag"`
	Sprint      string     `json:"sprint" db:"sprint"`
	Description string     `json:"description" db:"description"`
	Completed   bool       `json:"completed" db:"completed"`
	StartAt     *time.Time `json:"start_at" db:"start_at"`
	EndAt       *time.Time `json:"end_at" db:"end_at"`
}

func (t TaskEntity) GetId() int64 {
	return t.Id
}

func (t TaskEntity) GetTag() string {
	return t.Tag
}

func (t TaskEntity) GetSprint() string {
	return t.Sprint
}

func (t TaskEntity) GetDescription() string {
	return t.Description
}

func (t TaskEntity) GetCompleted() bool {
	return t.Completed
}

func (t TaskEntity) GetStartAt() *time.Time {
	return t.StartAt
}

func (t TaskEntity) GetEndAt() *time.Time {
	return t.EndAt
}

func (t *TaskEntity) SetId(id int64) m.Task {
	t.Id = id
	return t
}

func (t *TaskEntity) SetTag(tag string) m.Task {
	t.Tag = tag
	return t
}

func (t *TaskEntity) SetSprint(sprint string) m.Task {
	t.Sprint = sprint
	return t
}

func (t *TaskEntity) SetDescription(description string) m.Task {
	t.Description = description
	return t
}

func (t *TaskEntity) SetCompleted() m.Task {
	t.Completed = !t.Completed
	return t
}

func (t *TaskEntity) SetStartAt(startAt *time.Time) m.Task {
	t.StartAt = startAt
	return t
}

func (t *TaskEntity) SetEndAt(endAt *time.Time) m.Task {
	t.EndAt = endAt
	return t
}

func (t *TaskEntity) Localtime() m.Task {
	if t.GetStartAt() != nil {
		localtimeStartAt := t.GetStartAt().Local()
		t.SetStartAt(&localtimeStartAt)
	}
	if t.GetEndAt() != nil {
		localtimeEndAt := t.GetEndAt().Local()
		t.SetEndAt(&localtimeEndAt)
	}
	return t
}
