package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"host-lotus-worker/internal/domain"
	"host-lotus-worker/internal/repository"
)

type AssignmentProcessor struct {
	redisClient *redis.Client
	serverRepo *repository.ServerRepository
	worker      *domain.Worker
	interval    time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewAssignmentProcessor(redisClient *redis.Client, worker *domain.Worker, serverRepo *repository.ServerRepository, interval time.Duration) *AssignmentProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &AssignmentProcessor{
		redisClient: redisClient,
		serverRepo: serverRepo,
		worker:      worker,
		interval:    interval,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start begins the assignment processor background task
func (a *AssignmentProcessor) Start() {
	go a.run()
	log.Printf("Assignment processor started for worker %s with interval %v", a.worker.ID, a.interval)
}

// Stop gracefully stops the assignment processor background task
func (a *AssignmentProcessor) Stop() {
	a.cancel()
	log.Println("Assignment processor stopped")
}


func (a *AssignmentProcessor) run() {
    // Wait for heartbeat service to start and set worker state
    time.Sleep(a.interval)

	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	if err := a.checkAssignments(); err != nil {
		log.Printf("Failed to check initial assignments: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := a.checkAssignments(); err != nil {
				log.Printf("Failed to check assignments: %v", err)
			}
		case <-a.ctx.Done():
			log.Println("Assignment processor shutting down")
			return
		}
	}
}

// worker/internal/service/assignment_processor.go
const processAssignmentScript = `
    local key = KEYS[1]
    local serverID = ARGV[1]
    
    local data = redis.call('GET', key)
    if not data then
        return nil
    end
    
    local worker = cjson.decode(data)
    local assigned_servers = worker.assigned_servers or {}
    local running_servers = worker.running_servers or {}
    
    -- Check if this server is in assignments
    if not assigned_servers[serverID] then
        return redis.error_reply("Server not in assignments")
    end
    
    -- Move from assigned to running
    assigned_servers[serverID] = nil
    running_servers[serverID] = serverID
    
    worker.assigned_servers = assigned_servers
    worker.running_servers = running_servers
    
    redis.call('SET', key, cjson.encode(worker))
    return 1
`

func (a *AssignmentProcessor) processAssignment(serverID string) error {
    // TODO: actually start the server
    log.Printf("Starting server %s", serverID)

    server, err := a.serverRepo.GetServerById(a.ctx, serverID)
    if err != nil {
        return fmt.Errorf("failed to get server: %w", err)
    }

    server.Status = domain.StatusStarting
    err = a.serverRepo.UpdateServer(a.ctx, server)
    if err != nil {
        return fmt.Errorf("failed to update server: %w", err)
    }

    // Then atomically move from assigned to running
    a.worker.RunningServers[serverID] = serverID
    
    script := redis.NewScript(processAssignmentScript)
    return script.Run(a.ctx, a.redisClient,
        []string{fmt.Sprintf(workerKey, a.worker.ID)},
        serverID,
    ).Err()
}

func (a *AssignmentProcessor) checkAssignments() error {
    // Get current worker state from Redis
    key := fmt.Sprintf(workerKey, a.worker.ID)
    data, err := a.redisClient.Get(a.ctx, key).Result()
    if err != nil {
        return err
    }
    
    var worker domain.Worker
    json.Unmarshal([]byte(data), &worker)
    
    // Process any assigned servers
    for serverID := range worker.AssignedServers {
        // Start server and move to running
        err := a.processAssignment(serverID)
        if err != nil {
            log.Printf("Failed to process assignment: %v", err)
        }
    }

    return nil
}