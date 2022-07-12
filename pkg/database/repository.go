package database

import (
	"fmt"
	"sync"

	task "github.com/jorgefg4/todolist/pkg/task"
)

// Type to define a repository
type taskRepository struct {
	mtx   sync.RWMutex
	tasks map[int]*task.Task
	db    DatabaseHandler
}

// Returns a new repository initialized with the given tasks and
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

// Creates a new task using the database handler
func (r *taskRepository) CreateTask(g *task.Task) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.db.CreateNewTask(g.Name)

	return nil
}

// Retrieves the tasks in the database and updates the repository
// with them
func (r *taskRepository) FetchTasks() ([]*task.Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	tasks, err := r.db.GetAllTasks()

	r.tasks = tasks
	if err != nil {
		fmt.Println(err)
	}

	values := make([]*task.Task, 0, len(r.tasks))
	for _, value := range r.tasks {
		values = append(values, value)
	}

	return values, nil
}

// Deletes a task from the database
func (r *taskRepository) DeleteTask(ID int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.db.DeleteTask(ID)

	return nil
}

// Updates the status of a given task
func (r *taskRepository) UpdateTask(ID int) (int, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	err := r.db.CheckTask(ID)
	if err != nil {
		return 1, err
	} else {
		return 0, nil
	}
}

// Retrieves the status of a given task and updates the
// repository accordingly
func (r *taskRepository) FetchTaskByID(ID int) (*task.Task, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for _, v := range r.tasks {
		if v.ID == ID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("The ID %d doesn't exist", ID)
}

// Checks for the existance of a given task in the database
func (r *taskRepository) checkIfExists(ID int) error {
	for _, v := range r.tasks {
		if v.ID == ID {
			return fmt.Errorf("The task %d is already exist", ID)
		}
	}

	return nil
}
