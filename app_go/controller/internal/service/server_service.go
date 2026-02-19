package service

import (
	"context"
	"fmt"
	"host-lotus-controller/internal/domain"
	"host-lotus-controller/internal/repository"
	"time"
)

type ServerService struct {
	serverRepo *repository.ServerRepository
	workerRepo *repository.WorkerRepository
}

func NewServerService(serverRepo *repository.ServerRepository, workerRepo *repository.WorkerRepository) *ServerService {
	return &ServerService{
		serverRepo: serverRepo,
		workerRepo: workerRepo,
	}
}

func (s *ServerService) StartServer(ctx context.Context, user_id string, server_id string) error {
	// Check persistent server exists and verify ownership
	server_persistent, err := s.serverRepo.Persistent_repo.GetServerById(ctx, user_id, server_id)
	if err != nil {
		return err
	}

	if len(server_persistent) == 0 {
		return fmt.Errorf("server not found")
	}

	// Check if server is already running
	server_dynamic, err := s.serverRepo.Dynamic_repo.GetServerById(ctx, user_id, server_id)
	if err != nil {
		return err
	}

	if len(server_dynamic) != 0 && server_dynamic[0].Status != domain.StatusStopped {
		return fmt.Errorf("server is already running")
	}

	worker, err := s.workerRepo.GetFreeWorker(ctx)
	if err != nil {
		return err
	}

	// Add server to dynamic database
	err = s.serverRepo.AddServerDynamic(ctx, user_id, server_id, worker.PublicIp, domain.StatusScheduled, worker.ID)
	if err != nil {
		return err
	}

	// Assign server to worker
	err = s.workerRepo.AssignServerToWorker(ctx, worker.ID, server_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerService) StopServer(ctx context.Context, user_id string, server_id string) error {
	// Check persistent server exists and verify ownership
	server_persistent, err := s.serverRepo.Persistent_repo.GetServerById(ctx, user_id, server_id)
	if err != nil {
		return err
	}

	if len(server_persistent) == 0 {
		return fmt.Errorf("server not found")
	}

	// Check if server is running
	server_dynamic, err := s.serverRepo.Dynamic_repo.GetServerById(ctx, user_id, server_id)
	if err != nil {
		return err
	}

	if len(server_dynamic) == 0 || server_dynamic[0].Status != domain.StatusRunning {
		return fmt.Errorf("server is not running")
	}

	// Send command to worker
	err = s.workerRepo.SendCommandToWorker(ctx, server_dynamic[0].WorkerID, domain.ServerCommand{
		ServerID:  server_id,
		Action:    "stop",
		Timestamp: time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}
