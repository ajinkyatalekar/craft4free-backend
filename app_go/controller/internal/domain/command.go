package domain

import "time"

type ServerCommand struct {
    ServerID  string `json:"server_id"`
    Action    string `json:"action"`  // "start", "stop", "restart", "update"
	Timestamp time.Time `json:"timestamp"`
}