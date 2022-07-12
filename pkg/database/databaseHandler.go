package database

import (
	"github.com/jorgefg4/todolist/pkg/task"
)

// Interface for accesses to a database
type DatabaseHandler interface {
	// Stablishes connection with a given type of database
	GetConnection() error

	// Retrieves all tasks from a given type of database
	GetAllTasks() (map[int]*task.Task, error)

	// Creates a new task in a given type of database
	CreateNewTask(name string) error

	// Deletes a given task from a given type of database
	DeleteTask(id int) error

	// Modifies a given task from a given type of database
	CheckTask(id int) error
}
