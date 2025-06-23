package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	tasks     []task
	completed map[int]bool
	genextID  func() int
}

// addTask.
func addTask(currTask task, tasks []task) []task {
	fmt.Printf("Adding task: %d - %s\n", currTask.id, currTask.taskName)
	return append(tasks, currTask)
}

func ServeHomePage(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("Hello From Home")); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)

		return
	}
}

func (h *handler) ListAllTheTasks(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte(ListTasks(h.tasks, h.completed))); err != nil {
		log.Printf("Failed to List the Tasks: %v", err)
		http.Error(w, "Failed to List All the Tasks", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) ReturnSpecificID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	fmt.Println(parts[2])

	idStr := parts[2]
	taskID, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusNotFound)

		return
	}

	str := fmt.Sprintf("%+v", h.tasks[taskID])

	if _, err := w.Write([]byte(str)); err != nil {
		log.Printf("Internal Server Error %v", err)
		http.Error(w, "Failed to get that specific Task", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) CompleteThisTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusNotFound)

		return
	}

	markTaskAsCompleted(taskID, h.tasks, h.completed)

	if !h.completed[taskID] {
		http.Error(w, "Invalid task ID", http.StatusNotFound)

		return
	}

	if _, err := w.Write([]byte("SuccessFull Marked as Completed")); err != nil {
		log.Printf("Internal Server Error %v", err)
		http.Error(w, "Failed to Mark the Task as Completed", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) AddThisTask(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	var err error

	var body []byte

	// 1. Read the request body
	body, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Failed to read request body"}`, http.StatusBadRequest)

		return
	}

	fmt.Println(string(body))

	// 2. Define request structure
	var requestBody struct {
		TaskName  string `json:"taskName"`
		Completed bool   `json:"completed"`
	}

	// if You want to Marshal Any Slice of Bytes to a Struct , then You will have to keep
	// every field of that struct as Exported because :

	/*
		    The json.Marshal() function can only access exported (public) struct fields
			Your fields name and age are unexported (lowercase)
			By default, JSON marshaling silently ignores unexported fields
	*/

	// 3. Unmarshal using json.Unmarshal
	if err = json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)

		return
	}

	fmt.Println(requestBody)
	// 4. Validate required fields
	if requestBody.TaskName == "" {
		http.Error(w, `{"error": "taskName is required"}`, http.StatusBadRequest)

		return
	}

	// 5. Add the task
	newTask := task{
		id:        h.genextID(),
		taskName:  requestBody.TaskName,
		completed: requestBody.Completed,
	}
	h.tasks = append(h.tasks, newTask)

	// 6. Prepare response
	type res struct {
		ID        int    // `json :"id"
		TaskName  string // `json:"taskName"`
		Completed bool   // `json:"completed"`
		Message   string // `json:"message"`
	}

	response := res{
		ID:        newTask.id,
		TaskName:  newTask.taskName,
		Completed: newTask.completed,
		Message:   "Successfully Added Task",
	}

	// 7. Marshal and send response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate response"}`, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(jsonResponse); err != nil {
		log.Printf("Internal Server Error %v", err)
		http.Error(w, "Failed to Add the Task", http.StatusInternalServerError)

		return
	}
}

func main() {
	getNextID := generateID()
	h := &handler{
		tasks:     make([]task, 0),
		genextID:  getNextID,
		completed: make(map[int]bool),
	}

	currTask1 := task{
		getNextID(), "DemoTask1", false,
	}

	currTask2 := task{
		getNextID(), "DemoTask2", false,
	}
	currTask3 := task{
		getNextID(), "DemoTask3", false,
	}
	h.tasks = addTask(currTask1, h.tasks)
	h.tasks = addTask(currTask2, h.tasks)
	h.tasks = addTask(currTask3, h.tasks)

	http.HandleFunc("/", ServeHomePage)
	http.HandleFunc("GET /tasks", h.ListAllTheTasks) // geta all tasks
	http.HandleFunc("POST /tasks", h.AddThisTask)
	http.HandleFunc("GET /tasks/{id}", h.ReturnSpecificID) // get id
	http.HandleFunc("PUT /tasks/{id}", h.CompleteThisTask) // markComplete

	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
