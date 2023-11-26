package models

import "time"

type Task interface {
	GetId() int64
	GetTag() string
	GetDescription() string
	GetCompleted() bool
	GetStartAt() *time.Time
	GetEndAt() *time.Time

	SetId(id int64) Task
	SetTag(tag string) Task
	SetDescription(description string) Task
	SetCompleted() Task
	SetStartAt(startAt *time.Time) Task
	SetEndAt(endAt *time.Time) Task

	Localtime() Task
}
