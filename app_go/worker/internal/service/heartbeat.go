package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"host-lotus-worker/internal/domain"
)

const (
	workerKey = "worker:state:%s"
)

type HeartbeatService struct {
	redisClient *redis.Client
	worker      *domain.Worker
	interval    time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewHeartbeatService(redisClient *redis.Client, worker *domain.Worker, interval time.Duration) *HeartbeatService {
	ctx, cancel := context.WithCancel(context.Background())
	return &HeartbeatService{
		redisClient: redisClient,
		worker:      worker,
		interval:    interval,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start begins the heartbeat background task
func (h *HeartbeatService) Start() {
	go h.run()
	log.Printf("Heartbeat service started for worker %s with interval %v", h.worker.ID, h.interval)
}

// Stop gracefully stops the heartbeat background task
func (h *HeartbeatService) Stop() {
	h.cancel()
	log.Println("Heartbeat service stopped")
}

func (h *HeartbeatService) run() {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	if err := h.sendHeartbeat(); err != nil {
		log.Printf("Failed to send initial heartbeat: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := h.sendHeartbeat(); err != nil {
				log.Printf("Failed to send heartbeat: %v", err)
			}
		case <-h.ctx.Done():
			log.Println("Heartbeat service shutting down")
			return
		}
	}
}

const heartbeatScript = `
    local key = KEYS[1]
    local workerJSON = ARGV[1]
    
    local data = redis.call('GET', key)
    if not data then
        -- First heartbeat: initialize with full schema
        local worker = cjson.decode(workerJSON)
        worker.assigned_servers = {}  -- Empty array
        redis.call('SET', key, cjson.encode(worker))
        return 1
    end
    
    -- Existing worker: preserve controller-managed fields
    local existing = cjson.decode(data)
    local newWorker = cjson.decode(workerJSON)
    
    -- Preserve controller-managed fields
    newWorker.assigned_servers = existing.assigned_servers or {}
    
    redis.call('SET', key, cjson.encode(newWorker))
    return 1
`

func (h *HeartbeatService) sendHeartbeat() error {
    h.worker.LastHeartbeat = time.Now()
    
    workerData, _ := json.Marshal(h.worker)
    
    script := redis.NewScript(heartbeatScript)
    return script.Run(h.ctx, h.redisClient,
        []string{fmt.Sprintf(workerKey, h.worker.ID)},
        string(workerData),
    ).Err()
}