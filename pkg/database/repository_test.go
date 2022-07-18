package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	task "github.com/jorgefg4/todolist/pkg/task"
	"github.com/stretchr/testify/assert"
)

// TestCreateTask tests the creation of a new task
func TestCreateTask(t *testing.T) {
	var db *sql.DB
	var ctx context.Context = context.Background()

	postgreshandler := NewPostgres(db, ctx)

	// String to connect to database
	var conString string = "postgresql://" + os.Getenv("USER_DB") + ":" + os.Getenv("PASSWORD_DB") +
		"@" + os.Getenv("HOST_DB") + ":" + os.Getenv("PORT_DB") + "/" +
		os.Getenv("NAME_DB") + "?sslmode=disable"

	err := postgreshandler.GetConnection(conString)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("error")
	}

	tasks, err := postgreshandler.GetAllTasks()
	if err != nil {
		fmt.Println(err)
		t.Fatalf("error")
	}

	repo := NewTaskRepository(tasks, postgreshandler)

	t1 := task.Task{ID: 1, Name: "tarea de prueba"}
	err = repo.CreateTask(&t1)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("error")
	}

	repo1, _ := repo.FetchTasks()
	for _, value := range repo1 {
		if value.Name == "tarea de prueba" {
			assert.Equal(t, value.Name, t1.Name, "The two tasks should be the same.")
			err = repo.DeleteTask(value.ID) // after the test the task is deleted to avoid its persistence
			if err != nil {
				fmt.Println(err)
				t.Fatalf("error")
			}
		}
	}

}
