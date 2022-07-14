package server

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/jorgefg4/todolist/pkg/task"
)

func TestPrueba(t *testing.T) {
	// handler := api.New()

	// server := httptest.NewServer(handler)
	// defer server.Close()

	// e := httpexpect.New(t, server.URL)
	e := httpexpect.New(t, "http://localhost:8000")

	e.GET("/tasks").
		Expect().
		Status(http.StatusOK).JSON().Array().Empty()

	task := task.Task{
		Name:       "prueba",
		CheckValid: true,
	}

	e.POST("/tasks").WithJSON(task).
		Expect().
		Status(http.StatusCreated).JSON()

	obj := e.GET("/tasks").
		Expect().
		Status(http.StatusOK).JSON().Array().Element(0).Object()

	obj.Value("name").String().Equal("prueba")

	// e.DELETE("/tasks/100").
	// 	Expect().
	// 	Status(http.StatusNoContent).JSON()

	e.PUT("/tasks/153").
		Expect().
		Status(http.StatusBadRequest)

}
