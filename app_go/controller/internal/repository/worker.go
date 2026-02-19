package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"host-lotus-controller/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/supabase-community/supabase-go"
)

const (
	workerKey        = "worker:state:%s"
	allWorkersKey    = "worker:state"
	workerCommandKey = "worker:command:%s"
)

type WorkerRepository struct {
	redis_client    *redis.Client
	supabase_client *supabase.Client
}

func NewWorkerRepository(redis_client *redis.Client, supabase_client *supabase.Client) *WorkerRepository {
	return &WorkerRepository{
		redis_client:    redis_client,
		supabase_client: supabase_client,
	}
}

func unmarshalWorker(workerJson string) (domain.Worker, error) {
	var worker domain.Worker
	err := json.Unmarshal([]byte(workerJson), &worker)
	if err != nil {
		return domain.Worker{}, fmt.Errorf("failed to unmarshal worker data: %w", err)
	}

	return worker, nil
}

func (r *WorkerRepository) GetWorkerById(ctx context.Context, id string) (domain.Worker, error) {
	// Use Get to read the JSON string
	workerJson, err := r.redis_client.Get(ctx, fmt.Sprintf(workerKey, id)).Result()
	if err != nil {
		return domain.Worker{}, err
	}

	worker, err := unmarshalWorker(workerJson)
	if err != nil {
		return domain.Worker{}, err
	}

	return worker, nil
}

func (r *WorkerRepository) GetAllWorkers(ctx context.Context) ([]domain.Worker, error) {
	var workers []domain.Worker

	// Use SCAN to find all keys matching the pattern "worker:state:*"
	pattern := allWorkersKey + ":*"
	var cursor uint64
	var keys []string

	// Iterate through all keys matching the pattern
	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = r.redis_client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys: %w", err)
		}
		keys = append(keys, scanKeys...)

		// Break when cursor returns to 0 (iteration complete)
		if cursor == 0 {
			break
		}
	}

	// If no keys found, return empty slice
	if len(keys) == 0 {
		return workers, nil
	}

	// Get all worker data in a pipeline for efficiency
	pipe := r.redis_client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	// Unmarshal each worker
	for _, cmd := range cmds {
		workerJson, err := cmd.Result()
		if err != nil {
			continue
		}

		worker, err := unmarshalWorker(workerJson)
		if err != nil {
			return nil, fmt.Errorf("failed to get worker data: %w", err)
		}

		workers = append(workers, worker)
	}

	return workers, nil
}

func (r *WorkerRepository) GetOnlineWorkers(ctx context.Context) ([]domain.Worker, error) {
	allWorkers, err := r.GetAllWorkers(ctx)
	if err != nil {
		return nil, err
	}

	var onlineWorkers []domain.Worker
	for _, worker := range allWorkers {
		if worker.LastHeartbeat.After(time.Now().Add(-5 * time.Second)) {
			onlineWorkers = append(onlineWorkers, worker)
		}
	}

	return onlineWorkers, nil
}

func (r *WorkerRepository) GetFreeWorkers(ctx context.Context) ([]domain.Worker, error) {
	allWorkers, err := r.GetOnlineWorkers(ctx)
	if err != nil {
		return nil, err
	}

	var freeWorkers []domain.Worker
	for _, worker := range allWorkers {
		if len(worker.RunningServers)+len(worker.AssignedServers) < worker.MaxServers {
			freeWorkers = append(freeWorkers, worker)
		}
	}

	return freeWorkers, nil
}

func (r *WorkerRepository) GetFreeWorker(ctx context.Context) (domain.Worker, error) {
	workers, err := r.GetFreeWorkers(ctx)
	if err != nil {
		return domain.Worker{}, err
	}

	if len(workers) == 0 {
		return domain.Worker{}, fmt.Errorf("no free workers found")
	}

	return workers[0], nil
}

const assignServerScript = `
    local key = KEYS[1]
    local serverID = ARGV[1]
    local data = redis.call('GET', key)
    
    if not data then
        return redis.error_reply("Worker doesn't exist")
    end
    
    local worker = cjson.decode(data)
    local running_servers = worker.running_servers or {}
    local assigned_servers = worker.assigned_servers or {}
    local max_servers = worker.max_servers or 0
    
    -- Count servers in both maps
    local running_count = 0
    for _ in pairs(running_servers) do
        running_count = running_count + 1
    end
    
    local assigned_count = 0
    for _ in pairs(assigned_servers) do
        assigned_count = assigned_count + 1
    end
    
    -- Check capacity (running + pending assignments)
    if (running_count + assigned_count) >= max_servers then
        return redis.error_reply("Worker at capacity")
    end
    
    -- Add to assigned_servers map
    assigned_servers[serverID] = serverID
    worker.assigned_servers = assigned_servers
    redis.call('SET', key, cjson.encode(worker))
    
    return assigned_count + 1
`

// Assigns given server to a specific worker
func (r *WorkerRepository) AssignServerToWorker(ctx context.Context, workerID string, serverID string) error {
	script := redis.NewScript(assignServerScript)

	result, err := script.Run(ctx, r.redis_client,
		[]string{fmt.Sprintf(workerKey, workerID)},
		serverID,
	).Result()

	if err != nil {
		return err
	}

	if result == nil {
		return fmt.Errorf("failed to assign server to worker %s", workerID)
	}

	return nil
}

// Assigns given server to a free worker
func (r *WorkerRepository) AssignServer(ctx context.Context, serverID string) error {
	workers, err := r.GetFreeWorkers(ctx)
	if err != nil {
		return err
	}

	err = r.AssignServerToWorker(ctx, workers[0].ID, serverID)
	if err != nil {
		return err
	}

	return nil
}

func (r *WorkerRepository) SendCommandToWorker(ctx context.Context, workerID string, command domain.ServerCommand) error {
	data, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to marshal command data: %w", err)
	}
	err = r.redis_client.RPush(ctx, fmt.Sprintf(workerCommandKey, workerID), data).Err()
	if err != nil {
		return fmt.Errorf("failed to send command to worker: %w", err)
	}

	return nil
}
