package database

import (
	"sync"

	task "github.com/jorgefg4/todolist/pkg/task"
)

// Type to define a repository
type taskRepository struct {
	mtx   sync.RWMutex
	tasks map[int]*task.Task
	db    DatabaseHandler
}

// NewTaskRepository returns a new repository initialized with the given tasks and
// database handler
func NewTaskRepository(tasks map[int]*task.Task, DB DatabaseHandler) task.TaskRepository {
	if tasks == nil {
		tasks = make(map[int]*task.Task)
	}

	return &taskRepository{
		tasks: tasks,
		db:    DB,
	}
}

// CreateTask creates a new task using the database handler
func (r *taskRepository) CreateTask(g *task.Task) error {

	return r.db.CreateNewTask(g.Name)
}

// FetchTasks retrieves all the tasks in the database and updates
// the repository with them
func (r *taskRepository) FetchTasks() ([]*task.Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	tasks, err := r.db.GetAllTasks()

	r.tasks = tasks
	if err != nil {
		return nil, err
	}

	values := make([]*task.Task, 0, len(r.tasks))
	for _, value := range r.tasks {
		values = append(values, value)
	}

	return values, nil
}

// DeleteTask deletes a task from the database
func (r *taskRepository) DeleteTask(ID int) error {

	return r.db.DeleteTask(ID)
}

// CheckTask marks a given task as done
func (r *taskRepository) CheckTask(ID int) error {

	return r.db.CheckTask(ID)
}
