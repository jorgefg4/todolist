package database

import (
	"github.com/jorgefg4/todolist/pkg/task"
)

// Interface for accesses to a database
type DatabaseHandler interface {
	// GetConnection stablishes connection with a given type of database
	GetConnection(conString string) error

	// GetAllTasks retrieves all tasks from a given type of database
	GetAllTasks() (map[int]*task.Task, error)

	// CreateNewTask creates a new task in a given type of database
	CreateNewTask(name string) error

	// DeleteTask deletes a given task from a given type of database
	DeleteTask(id int) error

	// CheckTask modifies a given task from a given type of database
	CheckTask(id int) error
}
