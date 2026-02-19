package domain

import "time"

type Worker struct {
    // Presistent fields
    ID          string       `json:"id"`
	OciId       string       `json:"oci_id"`
	CreatedAt   time.Time    `json:"created_at"`
	PublicIp    string       `json:"public_ip"`
	PrivateIp   string       `json:"private_ip"`
	MaxServers  int          `json:"max_servers"`

	// Dynamic fields
	LastHeartbeat time.Time    	`json:"last_heartbeat"`
	RunningServers map[string]string     `json:"running_servers"`
	AssignedServers map[string]string    `json:"assigned_servers"`
}