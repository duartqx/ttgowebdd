package filters

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	f "github.com/duartqx/ttgowebdd/domains/filters"
)

const (
	TaskCompleted    = 1
	TaskNotCompleted = 2
)

type TaskFilter struct {
	Tag       string `json:"tag"`
	Completed int    `json:"completed"`
	Sprint    []int  `json:"sprint"`

	FromStartAt *time.Time `json:"from_start_at"`
	ToStartAt   *time.Time `json:"to_start_at"`

	FromEndAt *time.Time `json:"from_end_at"`
	ToEndAt   *time.Time `json:"to_end_at"`

	query  string
	values *[]interface{}
}

func NewTaskFilter() *TaskFilter {
	return &TaskFilter{
		values: &[]interface{}{},
	}
}

func (tf *TaskFilter) UnmarshalJSON(data []byte) error {
	var aux struct {
		Tag       []string `json:"tag"`
		Completed []string `json:"completed"`
		Sprint    []string `json:"sprint"`

		FromStartAt []string `json:"from_start_at"`
		ToStartAt   []string `json:"to_start_at"`

		FromEndAt []string `json:"from_end_at"`
		ToEndAt   []string `json:"to_end_at"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var sprints []int
	for _, s := range aux.Sprint {
		if si, err := strconv.Atoi(s); err == nil {
			sprints = append(sprints, si)
		}
	}
	tf.SetSprint(sprints)

	if len(aux.Tag) != 0 {
		tf.SetTag(aux.Tag[0])
	}

	if len(aux.Completed) != 0 {
		completed, err := strconv.Atoi(aux.Completed[0])
		if err != nil {
			return err
		}
		tf.SetCompleted(completed)
	}

	reflectSource := struct {
		FromStartAt string
		ToStartAt   string
		FromEndAt   string
		ToEndAt     string
	}{}

	if len(aux.FromStartAt) != 0 {
		reflectSource.FromStartAt = aux.FromStartAt[0]
	}
	if len(aux.ToStartAt) != 0 {
		reflectSource.ToStartAt = aux.ToStartAt[0]
	}
	if len(aux.FromEndAt) != 0 {
		reflectSource.FromEndAt = aux.FromEndAt[0]
	}
	if len(aux.ToEndAt) != 0 {
		reflectSource.ToEndAt = aux.ToEndAt[0]
	}

	format := "2006-01-02"
	target := reflect.ValueOf(tf).Elem()
	source := reflect.ValueOf(reflectSource)
	for _, field := range []string{"FromStartAt", "ToStartAt", "FromEndAt", "ToEndAt"} {
		value := source.FieldByName(field).String()
		if value != "" {
			timeParsed, err := time.Parse(format, value)
			if err != nil {
				return err
			}
			timeValue := reflect.ValueOf(&timeParsed)
			targetField := target.FieldByName(field)
			if targetField.IsValid() && targetField.CanSet() {
				targetField.Set(timeValue)
			}
		}
	}

	return nil
}

func (tf TaskFilter) GetTag() string {
	return tf.Tag
}

func (tf TaskFilter) GetSprints() []int {
	return tf.Sprint
}

func (tf TaskFilter) GetCompleted() int {
	return tf.Completed
}

func (tf TaskFilter) GetFromStartAt() *time.Time {
	return tf.FromStartAt
}

func (tf TaskFilter) GetToStartAt() *time.Time {
	return tf.ToStartAt
}

func (tf TaskFilter) GetFromEndAt() *time.Time {
	return tf.FromEndAt
}

func (tf TaskFilter) GetToEndAt() *time.Time {
	return tf.ToEndAt
}

func (tf *TaskFilter) SetTag(tag string) f.TaskFilter {
	tf.Tag = tag
	return tf
}

func (tf *TaskFilter) SetSprint(sprint []int) f.TaskFilter {
	tf.Sprint = sprint
	return tf
}

func (tf *TaskFilter) SetCompleted(completed int) f.TaskFilter {
	if slices.Contains[[]int](
		[]int{0, TaskCompleted, TaskNotCompleted}, completed,
	) {
		tf.Completed = completed
	}
	return tf
}

func (tf *TaskFilter) SetFromStartAt(from *time.Time) f.TaskFilter {
	tf.FromStartAt = from
	return tf
}

func (tf *TaskFilter) SetToStartAt(to *time.Time) f.TaskFilter {
	tf.ToStartAt = to
	return tf
}

func (tf *TaskFilter) SetFromEndAt(from *time.Time) f.TaskFilter {
	tf.FromEndAt = from
	return tf
}

func (tf *TaskFilter) SetToEndAt(to *time.Time) f.TaskFilter {
	tf.ToEndAt = to
	return tf
}

func (tf TaskFilter) GetCompletedQuery() string {
	switch tf.GetCompleted() {
	case 1:
		return "completed = 1"
	case 2:
		return "completed = 0"
	default:
		return ""
	}
}

func (tf TaskFilter) GetStartAtQuery() (string, *[]interface{}) {
	switch {
	case tf.GetFromStartAt() == nil && tf.GetToStartAt() == nil:
		// return "start_at BETWEEN ? AND ?", &[]interface{}{time.Now().AddDate(0, 0, -1), time.Now()}
		return "", &[]interface{}{}
	case tf.GetFromStartAt() == nil && tf.GetToStartAt() != nil:
		return "start_at <= ?", &[]interface{}{tf.GetToStartAt()}
	case tf.GetFromStartAt() != nil && tf.GetToStartAt() == nil:
		return "start_at >= ?", &[]interface{}{tf.GetFromStartAt()}
	case tf.GetFromStartAt() != nil && tf.GetToStartAt() != nil:
		return "start_at BETWEEN ? AND ?", &[]interface{}{tf.GetFromStartAt(), tf.GetToStartAt()}
	}

	return "", &[]interface{}{}
}

func (tf TaskFilter) GetEndAtQuery() (string, *[]interface{}) {

	switch {
	case tf.GetFromEndAt() == nil && tf.GetToEndAt() == nil:
		return "", &[]interface{}{}
	case tf.GetFromEndAt() == nil && tf.GetToEndAt() != nil:
		return "end_at <= ?", &[]interface{}{tf.GetToEndAt()}
	case tf.GetFromEndAt() != nil && tf.GetToEndAt() == nil:
		return "end_at >= ?", &[]interface{}{tf.GetFromEndAt()}
	case tf.GetFromEndAt() != nil && tf.GetToEndAt() != nil:
		return "end_at BETWEEN ? AND ?", &[]interface{}{tf.GetFromEndAt(), tf.GetToEndAt()}
	}

	return "", &[]interface{}{}
}

func (tf *TaskFilter) addToQuery(q string, values ...interface{}) {
	switch {
	case tf.query != "":
		tf.query += " AND " + q
	default:
		tf.query += q
	}
	*tf.values = append(*tf.values, values...)
}

func (tf *TaskFilter) Build() (string, *[]interface{}) {

	// Completed
	tf.query += tf.GetCompletedQuery()

	// StartAt
	startAtQuery, startAtValues := tf.GetStartAtQuery()
	if startAtQuery != "" {
		tf.addToQuery(startAtQuery, *startAtValues...)
	}

	// EndAt
	endAtQuery, endAtValues := tf.GetEndAtQuery()
	if endAtQuery != "" {
		tf.addToQuery(endAtQuery, *endAtValues...)
	}

	// Sprints
	sprintLen := len(tf.GetSprints())
	if sprintLen != 0 {
		var sprints []interface{}
		for _, s := range tf.GetSprints() {
			sprints = append(sprints, s)
		}
		var placeholder string
		for range sprintLen {
			placeholder += "?,"
		}
		placeholder, _ = strings.CutSuffix(placeholder, ",")
		tf.addToQuery(fmt.Sprintf("sprint IN (%s)", placeholder), sprints...)
	}

	if tf.query != "" {
		tf.query = " WHERE " + tf.query
	}

	return tf.query, tf.values
}
