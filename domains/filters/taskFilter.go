package filters

import "time"

type TaskFilter interface {
	GetTag() string
	GetCompleted() int
	GetSprints() []int
	GetFromStartAt() *time.Time
	GetToStartAt() *time.Time
	GetFromEndAt() *time.Time
	GetToEndAt() *time.Time

	SetTag(tag string) TaskFilter
	SetSprint(sprint []int) TaskFilter
	SetCompleted(completed int) TaskFilter
	SetFromStartAt(from *time.Time) TaskFilter
	SetToStartAt(to *time.Time) TaskFilter
	SetFromEndAt(from *time.Time) TaskFilter
	SetToEndAt(to *time.Time) TaskFilter

	GetCompletedQuery() string
	GetStartAtQuery() (string, *[]interface{})
	GetEndAtQuery() (string, *[]interface{})

	Build() (string, *[]interface{})
}
