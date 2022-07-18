package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/jorgefg4/todolist/pkg/task"
	"github.com/jorgefg4/todolist/web/assets"
)

// Templates
var navigationBarHTML string
var homepageTpl *template.Template

// Variable to assign IDs to tasks
var numID = 0

func init() {
	navigationBarHTML = assets.MustAssetString("web/templates/navigation_bar.html")
	homepageHTML := assets.MustAssetString("web/templates/index.html")
	homepageTpl = template.Must(template.New("homepage_view").Parse(homepageHTML))
}

type api struct {
	router     http.Handler
	repository task.TaskRepository
}

type Server interface {
	Router() http.Handler
	fetchTasks(w http.ResponseWriter, r *http.Request) //for the test
}

// Router returns the router of the api (its a method)
func (a *api) Router() http.Handler {
	return a.router
}

// New returns a Server
func New(repo task.TaskRepository) Server {
	a := &api{repository: repo}

	r := mux.NewRouter() //router instance creation

	// swagger documentation for share
	r.Handle("/swagger.yml", http.FileServer(http.Dir("./")))
	opts1 := middleware.RedocOpts{Path: "api-doc", SpecURL: "/swagger.yml"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/api-doc", sh1)

	//endpoints:
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/tasks", a.fetchTasks).Methods(http.MethodGet)
	r.HandleFunc("/tasks", a.addTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{ID:[a-zA-Z0-9_]+}", a.removeTask).Methods(http.MethodDelete)
	r.HandleFunc("/tasks/{ID:[a-zA-Z0-9_]+}", a.modifyTask).Methods(http.MethodPut)
	r.PathPrefix("/web/static/").Handler(http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static"))))

	a.router = r
	return a
}

// HomeHandler renders the homepage view template
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	push(w, "/web/static/style.css")
	push(w, "/web/static/todolist.css")
	push(w, "/web/static/navigation_bar.css")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fullData := map[string]interface{}{
		"NavigationBar": template.HTML(navigationBarHTML),
	}
	render(w, r, homepageTpl, "homepage_view", fullData)
}

// push pushes the given resource to the client
func push(w http.ResponseWriter, resource string) {
	pusher, ok := w.(http.Pusher)

	if ok {
		err := pusher.Push(resource, nil)
		log.Fatal(err)
	}
}

// render renders a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	_, err := w.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

///////// Handlers: /////////////

// fetchTasks shows all the tasks
func (a *api) fetchTasks(w http.ResponseWriter, r *http.Request) {
	tasks, _ := a.repository.FetchTasks()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Fatal(err)
	}
}

// addTask adds a new task
func (a *api) addTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var t task.Task
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")

	numID++ //first the ID is incremented
	t.ID = numID
	err = a.repository.CreateTask(&t)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(numID) //the ID of the task is sent as a response
	if err != nil {
		log.Fatal(err)
	}
}

// removeTask removes a existing task
func (a *api) removeTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["ID"]) //string to int ID conversion

	err := a.repository.DeleteTask(id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// modifyTask used to mark a task as done
func (a *api) modifyTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["ID"]) //string to int ID conversion

	err := a.repository.CheckTask(id)
	if err != nil { //if an error is received, BadRequest 404 is displayed (the indicated task does not exist)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
