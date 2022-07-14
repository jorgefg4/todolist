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

//Variable para asignar IDs a las tasks
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
	fetchTasks(w http.ResponseWriter, r *http.Request) //para el test
}

func New(repo task.TaskRepository) Server {
	a := &api{repository: repo}

	r := mux.NewRouter() //creamos una instancia de router

	//r.Handle("/swagger.yml", http.FileServer(http.Dir("./")))

	// documentation for share
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
	//r.HandleFunc("/tasks/{ID:[a-zA-Z0-9_]+}", a.fetchTask).Methods(http.MethodGet)

	a.router = r
	return a
}

func (a *api) Router() http.Handler { //metodo
	return a.router
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

// Push the given resource to the client.
func push(w http.ResponseWriter, resource string) {
	pusher, ok := w.(http.Pusher)

	if ok {
		err := pusher.Push(resource, nil)
		log.Fatal(err)
	}
}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

//Handlers:
func (a *api) fetchTasks(w http.ResponseWriter, r *http.Request) { //para mostrar todas las tareas
	tasks, _ := a.repository.FetchTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (a *api) addTask(w http.ResponseWriter, r *http.Request) { //para añadir nueva tarea
	decoder := json.NewDecoder(r.Body)

	var t task.Task
	decoder.Decode(&t)

	w.Header().Set("Content-Type", "application/json")

	numID++ //primero incrementamos el ID
	t.ID = numID
	a.repository.CreateTask(&t)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(numID) //se envia como respuesta el ID de la task
}

func (a *api) removeTask(w http.ResponseWriter, r *http.Request) { //para borrar una tarea
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["ID"]) //convierto el ID de string a int
	a.repository.DeleteTask(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

//modifyTask sirve para marcar una tarea como realizada
func (a *api) modifyTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["ID"]) //convierto el ID de string a int
	err := a.repository.CheckTask(id)

	if err != nil { //si se recibe error se muestra BadRequest 404 (la tarea indicada no existe)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}
