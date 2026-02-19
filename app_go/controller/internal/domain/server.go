package domain

import "time"

type Server struct {
    // Presistent fields
    ID          string       `json:"id"`
    UserId      string       `json:"user_id"`
    Name        string       `json:"name"`
    Type        string       `json:"type"`
    Version     string       `json:"version"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    DeletedAt   *time.Time   `json:"deleted_at"`

    // Dynamic ONLY fields
    IP          string       `json:"ip"`
    Status      ServerStatus `json:"status"`
    WorkerID    string       `json:"worker_id"`
}

type ServerPersistent struct {
    // Presistent fields
    ID          string       `json:"id"`
    UserId      string       `json:"user_id"`
    Name        string       `json:"name"`
    Type        string       `json:"type"`
    Version     string       `json:"version"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    DeletedAt   *time.Time   `json:"deleted_at"`
}

type ServerDynamic struct {
    // Presistent fields
    ID          string       `json:"id"`
    UserId      string       `json:"user_id"`
    Name        string       `json:"name"`
    Type        string       `json:"type"`
    Version     string       `json:"version"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    DeletedAt   *time.Time   `json:"deleted_at"`

    // Dynamic fields
    IP          string       `json:"ip"`
    Status      ServerStatus `json:"status"`
    WorkerID    string       `json:"worker_id"`
}

type ServerStatus string
const (
    StatusScheduled ServerStatus = "scheduled"
    StatusRunning ServerStatus = "running"
    StatusStopped ServerStatus = "stopped"
    StatusStarting ServerStatus = "starting"
)