package database

import (
	"context"
	"database/sql"

	task "github.com/jorgefg4/todolist/pkg/task"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jorgefg4/todolist/models"

	_ "github.com/lib/pq"
)

// Type to access a Postgresql database
// Implements DatabaseHandler
type PostgresHandler struct {
	DB  *sql.DB
	ctx context.Context
}

// NewPostgres returns a new PostgresHandler element
func NewPostgres(db *sql.DB, ctx context.Context) *PostgresHandler {
	return &PostgresHandler{
		DB:  db,
		ctx: ctx,
	}
}

// GetConnection stablish a connection with the database
func (handler *PostgresHandler) GetConnection(conString string) error {
	var err error
	handler.DB, err = sql.Open("postgres", conString)
	if err != nil {
		return err
	}

	return handler.DB.Ping()
}

// GetAllTasks retrieves all tasks from the database and returns a map with the tasks formatted
// for the application to use (type of Task)
func (handler *PostgresHandler) GetAllTasks() (map[int]*task.Task, error) {
	tasks, err := models.Tasks().All(handler.ctx, handler.DB)
	if err != nil {
		return nil, err
	}

	tasksMap := make(map[int]*task.Task)

	for _, databaseTask := range tasks {
		id := databaseTask.ID
		name := databaseTask.Name
		check_valid := databaseTask.CheckValid

		newTask := task.Task{
			ID:         id,
			Name:       name,
			CheckValid: check_valid,
		}

		tasksMap[id] = &newTask
	}

	return tasksMap, nil
}

// CreateNewTask creates a new task in the database
func (handler *PostgresHandler) CreateNewTask(name string) error {
	var newTask models.Task

	newTask.Name = name

	return newTask.Insert(handler.ctx, handler.DB, boil.Infer())
}

// DeleteTask deletes a given task from the database
func (handler *PostgresHandler) DeleteTask(id int) error {
	task, err := models.FindTask(handler.ctx, handler.DB, id)
	if err != nil {
		return err
	}

	_, err = task.Delete(handler.ctx, handler.DB)

	return err
}

// CheckTask marks a given task as done
func (handler *PostgresHandler) CheckTask(id int) error {
	task, err := models.FindTask(handler.ctx, handler.DB, id)
	if err != nil {
		return err
	}

	task.CheckValid = true

	_, err = task.Update(handler.ctx, handler.DB, boil.Infer())

	return err
}
