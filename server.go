package main

import (
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var (
	router  *mux.Router
	session *r.Session
)

func init() {
	var err error

	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "todo",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func NewServer(addr string) *http.Server {
	// Setup router
	router = initRouting()

	// Create and start server
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func StartServer(server *http.Server) {
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error: %v", err)
	}
}

func initRouting() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/active", activeIndexHandler)
	r.HandleFunc("/completed", completedIndexHandler)
	r.HandleFunc("/new", newHandler)
	r.HandleFunc("/toggle/{id}", toggleHandler)
	r.HandleFunc("/delete/{id}", deleteHandler)
	r.HandleFunc("/clear", clearHandler)

	// Add handler for static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return r
}

// Handlers

func indexHandler(w http.ResponseWriter, req *http.Request) {
	items := []TodoItem{}

	// Fetch all the items from the database
	res, err := r.Table("items").OrderBy(r.Asc("Created")).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = res.All(&items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "index", items)
}

func activeIndexHandler(w http.ResponseWriter, req *http.Request) {
	items := []TodoItem{}

	// Fetch all the items from the database
	query := r.Table("items").Filter(r.Row.Field("Status").Eq("active"))
	query = query.OrderBy(r.Asc("Created"))
	res, err := query.Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = res.All(&items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "index", items)
}

func completedIndexHandler(w http.ResponseWriter, req *http.Request) {
	items := []TodoItem{}

	// Fetch all the items from the database
	query := r.Table("items").Filter(r.Row.Field("Status").Eq("complete"))
	query = query.OrderBy(r.Asc("Created"))
	res, err := query.Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = res.All(&items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "index", items)
}

func newHandler(w http.ResponseWriter, req *http.Request) {
	// Create the item
	item := NewTodoItem(req.PostFormValue("text"))
	item.Created = time.Now()

	// Insert the new item into the database
	_, err := r.Table("items").Insert(item).RunWrite(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

func toggleHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if id == "" {
		http.NotFound(w, req)
		return
	}

	// Check that the item exists
	res, err := r.Table("items").Get(id).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.IsNil() {
		http.NotFound(w, req)
		return
	}

	// Toggle the item
	_, err = r.Table("items").Get(id).Update(map[string]interface{}{"Status": r.Branch(
		r.Row.Field("Status").Eq("active"),
		"complete",
		"active",
	)}).RunWrite(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if id == "" {
		http.NotFound(w, req)
		return
	}

	// Check that the item exists
	res, err := r.Table("items").Get(id).Run(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.IsNil() {
		http.NotFound(w, req)
		return
	}

	// Delete the item
	_, err = r.Table("items").Get(id).Delete().RunWrite(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

func clearHandler(w http.ResponseWriter, req *http.Request) {
	// Delete all completed items
	_, err := r.Table("items").Filter(
		r.Row.Field("Status").Eq("complete"),
	).Delete().RunWrite(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)
}
