package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jorgefg4/todolist/pkg/database"
	"github.com/jorgefg4/todolist/pkg/task"
)

// TestFetchTasks tests the retrieval of the tasks from the database
func TestFetchTasks(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}

	h := database.PostgresHandler{}

	// String to connect to database
	var conString string = "postgresql://" + os.Getenv("USER_DB") + ":" + os.Getenv("PASSWORD_DB") +
		"@" + os.Getenv("HOST_DB") + ":" + os.Getenv("PORT_DB") + "/" +
		os.Getenv("NAME_DB") + "?sslmode=disable"

	err = h.GetConnection(conString)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("error")
	}

	tasks, err := h.GetAllTasks()
	if err != nil {
		fmt.Println(err)
		t.Fatalf("error")
	}

	repo := database.NewTaskRepository(tasks, &h)
	s := New(repo)

	rec := httptest.NewRecorder() // With the httptest packet its possible to generate the http.ResponseWriter

	s.fetchTasks(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got: %d", http.StatusOK, res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var got []task.Task
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("could not unmarshall response %v", err)
	}

	expected := len(tasks)

	if len(got) != expected {
		t.Errorf("expected %v tasks, got: %v task", expected, len(got))
	}
}
