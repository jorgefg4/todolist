package task

// Type that defines a task
// Tags are included to transform the struct from json or to json
type Task struct {
	ID         int    `json:"ID"`
	Name       string `json:"name,omitempty"`
	CheckValid bool   `json:"check_valid,omitempty"`
}

// Defines the interface to interact with the tasks
// Repository provides access to the task storage
type TaskRepository interface {

	// Creates a new task
	CreateTask(g *Task) error

	// Fetch the tasks in the database
	FetchTasks() ([]*Task, error)

	// Deletes a given task
	DeleteTask(ID int) error

	// Updates a given task
	CheckTask(ID int) error
}
