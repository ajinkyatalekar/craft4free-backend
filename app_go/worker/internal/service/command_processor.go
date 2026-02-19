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

const (
	workerCommandKey = "worker:command:%s"
)

// Action constants
const (
	ActionStop    = "stop"
	ActionRestart = "restart"
	ActionUpdate  = "update"
)

type ServerCommand struct {
	ServerID  string            `json:"server_id"`
	Action    string            `json:"action"`
	Params    map[string]string `json:"params,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

type CommandProcessor struct {
	redisClient *redis.Client
	serverRepo  *repository.ServerRepository
	worker      *domain.Worker
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewCommandProcessor(redisClient *redis.Client, worker *domain.Worker, serverRepo *repository.ServerRepository) *CommandProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &CommandProcessor{
		redisClient: redisClient,
		serverRepo:  serverRepo,
		worker:      worker,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (c *CommandProcessor) Start() {
	go c.run()
	log.Printf("Command processor started for worker %s", c.worker.ID)
}

func (c *CommandProcessor) Stop() {
	c.cancel()
	log.Println("Command processor stopped")
}

func (c *CommandProcessor) run() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// BLPOP blocks until a command is available
			result, err := c.redisClient.BLPop(c.ctx, 5*time.Second,
				fmt.Sprintf(workerCommandKey, c.worker.ID)).Result()
			if err == redis.Nil {
				continue
			}
			if err != nil {
				log.Printf("Error reading commands: %v", err)
				continue
			}

			var cmd ServerCommand
			json.Unmarshal([]byte(result[1]), &cmd)
			if err := c.processCommand(cmd); err != nil {
				log.Printf("Error processing command: %v", err)
			}
		}
	}
}

func (c *CommandProcessor) processCommand(cmd ServerCommand) error {
	switch cmd.Action {
	case ActionStop:
		return c.stopServer(cmd.ServerID)
	case ActionRestart:
		return c.restartServer(cmd.ServerID)
	case ActionUpdate:
		return c.updateServerConfig(cmd.ServerID, cmd.Params)
	default:
		return fmt.Errorf("unknown action: %s", cmd.Action)
	}
}

func (c *CommandProcessor) stopServer(serverID string) error {
	log.Printf("Stopping server %s", serverID)
	// TODO: implement actual server stop logic
	server, err := c.serverRepo.GetServerById(c.ctx, serverID)
	if err != nil {
		return err
	}
	server.Status = domain.StatusStopped
	err = c.serverRepo.UpdateServer(c.ctx, server)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommandProcessor) restartServer(serverID string) error {
	log.Printf("Restarting server %s", serverID)
	// TODO: implement actual server restart logic
	return nil
}

func (c *CommandProcessor) updateServerConfig(serverID string, params map[string]string) error {
	log.Printf("Updating server %s config with params: %v", serverID, params)
	// TODO: implement actual config update logic
	return nil
}
