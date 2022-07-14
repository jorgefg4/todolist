package service

import (
	"fmt"

	"github.com/jorgefg4/todolist/pkg/database"
	"github.com/jorgefg4/todolist/pkg/server"
)

// Type to define the service containing the access to a database
type Service struct {
	DB database.DatabaseHandler
}

// Returns a new service with the given database handler
func NewService(DB database.DatabaseHandler) *Service {
	return &Service{
		DB: DB,
	}
}

// Returns a new server that connects to a database
func (svc *Service) NewServer(conString string) (server.Server, error) {
	err := svc.DB.GetConnection(conString)
	if err != nil {
		return nil, err
	}

	tasks, err := svc.DB.GetAllTasks()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		fmt.Println(task.Name)
	}

	// Llamo al package "server" para crear un nuevo router
	repo := database.NewTaskRepository(tasks, svc.DB)
	s := server.New(repo)

	return s, err
}
