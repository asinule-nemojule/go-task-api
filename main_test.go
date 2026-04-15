package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestCreateAndGetTask(t *testing.T) {
	// Resetuj tasks pre testa
	tasksMu.Lock()
	tasks = make(map[string]Task)
	tasksMu.Unlock()

	// 1. Kreiraj task
	taskJSON := `{"title":"Test Task"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", w.Code)
	}

	var createdTask Task
	json.NewDecoder(w.Body).Decode(&createdTask)

	if createdTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", createdTask.Title)
	}

	// 2. Dohvati task
	req2 := httptest.NewRequest("GET", "/tasks/"+createdTask.ID, nil)
	req2.SetPathValue("id", createdTask.ID)
	w2 := httptest.NewRecorder()

	getTask(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w2.Code)
	}

	var fetchedTask Task
	json.NewDecoder(w2.Body).Decode(&fetchedTask)

	if fetchedTask.ID != createdTask.ID {
		t.Errorf("Task ID mismatch")
	}
}

func TestDeleteTask(t *testing.T) {
	// Reset
	tasksMu.Lock()
	tasks = make(map[string]Task)
	tasksMu.Unlock()

	// Prvo kreiraj task
	taskJSON := `{"title":"To Be Deleted"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	createTask(w, req)

	var task Task
	json.NewDecoder(w.Body).Decode(&task)

	// Onda ga obriši
	req2 := httptest.NewRequest("DELETE", "/tasks/"+task.ID, nil)
	req2.SetPathValue("id", task.ID)
	w2 := httptest.NewRecorder()
	deleteTask(w2, req2)

	if w2.Code != http.StatusNoContent {
		t.Errorf("Expected NoContent, got %v", w2.Code)
	}

	// Proveri da više ne postoji
	req3 := httptest.NewRequest("GET", "/tasks/"+task.ID, nil)
	req3.SetPathValue("id", task.ID)
	w3 := httptest.NewRecorder()
	getTask(w3, req3)

	if w3.Code != http.StatusNotFound {
		t.Errorf("Expected NotFound, got %v", w3.Code)
	}
}

func TestGetAllTasks(t *testing.T) {
	// Reset
	tasksMu.Lock()
	tasks = make(map[string]Task)
	tasksMu.Unlock()

	// Dodaj dva taska
	task1 := Task{ID: "1", Title: "Task 1", Completed: false}
	task2 := Task{ID: "2", Title: "Task 2", Completed: true}
	tasks["1"] = task1
	tasks["2"] = task2

	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	getTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected OK, got %v", w.Code)
	}

	var taskList []Task
	json.NewDecoder(w.Body).Decode(&taskList)

	if len(taskList) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(taskList))
	}
}
