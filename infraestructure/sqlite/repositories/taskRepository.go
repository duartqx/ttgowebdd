package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	t "github.com/duartqx/ttgowebdd/domains/entities/task"
	m "github.com/duartqx/ttgowebdd/domains/models"

	f "github.com/duartqx/ttgowebdd/infraestructure/sqlite/filters"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tag TEXT NOT NULL,
			description TEXT NOT NULL,
			completed BOOL DEFAULT 0,
			start_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			end_at DATETIME DEFAULT NULL
		)
	`)
	return &TaskRepository{db: db}
}

func (tr TaskRepository) getModel() *t.TaskEntity {
	return &t.TaskEntity{}
}

func (tr TaskRepository) Filter(reader io.ReadCloser) (*[]m.Task, error) {

	tf := f.NewTaskFilter()
	if reader != nil {
		if err := json.NewDecoder(reader).Decode(&tf); err != nil {
			return nil, fmt.Errorf("Decode Error: %v", err)
		}
	}

	if tf.GetTag() != "" {
		task, err := tr.FindByTag(tf.GetTag())
		return &[]m.Task{task}, err
	}

	where, whereValues := tf.Build()

	rows, err := tr.db.Query(
		"SELECT id, tag, description, completed, start_at, end_at FROM tasks"+where,
		*whereValues...,
	)
	if err != nil {
		return nil, err
	}

	tasks := []m.Task{}

	for rows.Next() {
		task := tr.getModel()

		if err := rows.Scan(
			&task.Id, &task.Tag, &task.Description, &task.Completed, &task.StartAt, &task.EndAt,
		); err != nil {
			return nil, err
		}

		var iTask m.Task = task
		tasks = append(tasks, iTask)
	}

	return &tasks, nil
}

func (tr TaskRepository) FindById(id int64) (m.Task, error) {
	task := tr.getModel()
	if err := tr.db.Get(task, "SELECT * FROM tasks WHERE id = ? LIMIT 1", id); err != nil {
		return nil, err
	}
	return task, nil
}

func (tr TaskRepository) FindByTag(tag string) (m.Task, error) {
	task := tr.getModel()
	if err := tr.db.Get(task, "SELECT * FROM tasks WHERE tag = ? LIMIT 1", tag); err != nil {
		return nil, err
	}
	return task, nil
}

func (tr TaskRepository) Create(task m.Task) error {
	var (
		taskId  int64
		startAt time.Time
	)
	if err := tr.db.QueryRow(
		"INSERT INTO tasks (tag, description) VALUES (?, ?) RETURNING id, start_at",
		task.GetTag(), task.GetDescription(),
	).Scan(&taskId, &startAt); err != nil {
		return err
	}

	task.SetId(taskId).SetStartAt(&startAt)

	return nil
}

func (tr TaskRepository) Complete(task m.Task) error {

	if task.GetId() == 0 {
		return fmt.Errorf("Invalid Error")
	}

	var endAt time.Time

	if err := tr.db.Get(
		&endAt,
		`UPDATE tasks SET
		completed = CASE WHEN completed = 1 THEN 0 ELSE 1 END,
		end_at = CASE WHEN completed = 1 THEN NULL ELSE CURRENT_TIMESTAMP END
		WHERE id = ?`,
		task.GetId(),
	); err != nil {
		return err
	}

	task.SetEndAt(&endAt).SetCompleted()

	return nil
}

func (tr TaskRepository) CompleteById(id int) (m.Task, error) {

	task := tr.getModel()

	if err := tr.db.Get(
		task,
		`UPDATE tasks SET
		completed = CASE WHEN completed = 1 THEN 0 ELSE 1 END,
		end_at = CASE WHEN completed = 1 THEN NULL ELSE CURRENT_TIMESTAMP END
		WHERE id = ?
		RETURNING *`,
		id,
	); err != nil {
		return nil, err
	}
	return task, nil
}

func (tr TaskRepository) UpdateEndAt(task m.Task) error {
	endAt := time.Now()
	if _, err := tr.db.Exec(
		"UPDATE tasks SET end_at = ? WHERE id = ?", endAt, task.GetId(),
	); err != nil {
		return err
	}

	task.SetEndAt(&endAt)

	return nil
}

func (tr TaskRepository) UpdateEndAtById(id int) (m.Task, error) {
	task := tr.getModel()
	if err := tr.db.Get(
		task,
		"UPDATE tasks SET end_at = ? WHERE id = ? RETURNING *", time.Now(), id,
	); err != nil {
		return nil, err
	}
	return task, nil
}
