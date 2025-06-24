// handler_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestHandler() *handler {
	getNextID := generateID()

	h := &handler{
		tasks:     make([]task, 0),
		genextID:  getNextID,
		completed: make(map[int]bool),
	}

	// adding Dummy Tasks for Testing
	h.tasks = append(h.tasks, task{id: getNextID(), taskName: "DemoTask 1", completed: false})
	h.tasks = append(h.tasks, task{id: getNextID(), taskName: "DemoTask 2", completed: false})

	return h
}

func TestListAllTheTasks(t *testing.T) {
	h := setupTestHandler()
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	h.ListAllTheTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Check if it contains our test tasks
	body := rr.Body.String()
	if !strings.Contains(body, "DemoTask 1") || !strings.Contains(body, "DemoTask 2") {
		t.Error("response should contain  All the Demo tasks")
	}
}
func TestReturnSpecificID(t *testing.T) {
	h := setupTestHandler()

	tests := []struct {
		name       string
		taskID     string
		wantStatus int
	}{
		{"Valid ID", "1", http.StatusOK},
		{"Invalid ID", "999", http.StatusNotFound},
		{"Non-numeric ID", "abc", http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/tasks/{id}", nil)
		req.SetPathValue("id", tt.taskID)
		rr := httptest.NewRecorder()

		h.ReturnSpecificID(rr, req)

		if rr.Code != tt.wantStatus {
			t.Errorf("expected status %d, got %d", tt.wantStatus, rr.Code)
		}

		// if its Okay then Please check if Body is Same Expected or Not
		if tt.wantStatus == http.StatusOK && !strings.Contains(rr.Body.String(), "DemoTask 1") {
			t.Error("response should contain the task")
		}
	}
}

func TestCompleteThisTask(t *testing.T) {
	h := setupTestHandler()

	tests := []struct {
		name       string
		taskID     string
		wantStatus int
	}{
		{"Valid ID", "1", http.StatusOK},
		{"Invalid ID", "999", http.StatusNotFound},
		{"Non-numeric ID", "abc", http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("PUT", "/tasks/{id}", nil)
		req.SetPathValue("id", tt.taskID)
		rr := httptest.NewRecorder()

		h.CompleteThisTask(rr, req)

		if rr.Code != tt.wantStatus {
			t.Errorf("expected status %d, got %d", tt.wantStatus, rr.Code)
		}

		if tt.wantStatus == rr.Code && !h.completed[1] {
			t.Error("task should be marked as completed")
		}
	}
}

func TestAddThisTask(t *testing.T) {
	h := setupTestHandler()

	tests := []struct {
		name        string
		requestBody map[string]interface{}
		wantStatus  int
	}{
		{
			"Valid Task",
			map[string]interface{}{"taskName": "New Task", "completed": false},
			http.StatusCreated,
		},
		{
			"Missing Task Name",
			map[string]interface{}{"completed": false},
			http.StatusBadRequest,
		},
		{
			"Invalid JSON",
			nil,
			http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		var body bytes.Buffer
		if tt.requestBody != nil {
			json.NewEncoder(&body).Encode(tt.requestBody)
		}

		req := httptest.NewRequest("POST", "/tasks", &body)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		h.AddThisTask(rr, req)

		if rr.Code != tt.wantStatus {
			t.Errorf("expected status %d, got %d", tt.wantStatus, rr.Code)
		}

		if tt.wantStatus == http.StatusCreated {
			var response struct {
				ID        int    `json:"ID"`
				TaskName  string `json:"TaskName"`
				Completed bool   `json:"Completed"`
			}
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Fatal("failed to decode response")
			}

			if response.TaskName != "New Task" {
				t.Errorf("expected task name 'New Task', got '%s'", response.TaskName)
			}

			if len(h.tasks) != 3 {
				t.Errorf("expected 3 tasks, got %d", len(h.tasks))
			}
		}
	}
}
