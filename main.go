package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
)

// Task struktura - naš model podataka
type Task struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
}

// In-memory "baza" podataka
var (
    tasks   = make(map[string]Task)
    tasksMu sync.RWMutex
)

func main() {
    // Rute (endpoint-i)
    http.HandleFunc("GET /tasks", getTasks)
    http.HandleFunc("GET /tasks/{id}", getTask)
    http.HandleFunc("POST /tasks", createTask)
    http.HandleFunc("DELETE /tasks/{id}", deleteTask)
    http.HandleFunc("GET /health", healthCheck)

    port := ":8080"
    fmt.Printf("Server running on http://localhost%s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

// GET /tasks - vraća sve taskove
func getTasks(w http.ResponseWriter, r *http.Request) {
    tasksMu.RLock()
    defer tasksMu.RUnlock()

    taskList := make([]Task, 0, len(tasks))
    for _, task := range tasks {
        taskList = append(taskList, task)
    }

    sendJSON(w, http.StatusOK, taskList)
}

// GET /tasks/{id} - vraća jedan task
func getTask(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    tasksMu.RLock()
    defer tasksMu.RUnlock()

    task, exists := tasks[id]
    if !exists {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    sendJSON(w, http.StatusOK, task)
}

// POST /tasks - kreira novi task
func createTask(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Title string `json:"title"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }

    tasksMu.Lock()
    defer tasksMu.Unlock()

    id := fmt.Sprintf("%d", time.Now().UnixNano())
    task := Task{
        ID:        id,
        Title:     req.Title,
        Completed: false,
        CreatedAt: time.Now(),
    }
    tasks[id] = task

    sendJSON(w, http.StatusCreated, task)
}

// DELETE /tasks/{id} - briše task
func deleteTask(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    tasksMu.Lock()
    defer tasksMu.Unlock()

    if _, exists := tasks[id]; !exists {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    delete(tasks, id)
    w.WriteHeader(http.StatusNoContent)
}

// GET /health - health check za Docker i monitoring
func healthCheck(w http.ResponseWriter, r *http.Request) {
    sendJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Helper funkcija za slanje JSON odgovora
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}