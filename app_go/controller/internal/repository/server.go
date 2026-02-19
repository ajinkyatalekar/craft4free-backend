package repository

import (
	"context"
	"fmt"
	"host-lotus-controller/internal/domain"

	"github.com/redis/go-redis/v9"
	"github.com/supabase-community/supabase-go"
)

type ServerRepository struct {
	supabase_client *supabase.Client
	redis_client    *redis.Client
	Dynamic_repo    *serverDynamicRepository
	Persistent_repo *serverPersistentRepository
}

func NewServerRepository(redis_client *redis.Client, supabase_client *supabase.Client) *ServerRepository {
	return &ServerRepository{
		supabase_client: supabase_client,
		redis_client:    redis_client,
		Dynamic_repo:    newServerDynamicRepository(supabase_client, redis_client),
		Persistent_repo: newServerPersistentRepository(supabase_client, redis_client),
	}
}

func (r *ServerRepository) GetAllServers(ctx context.Context, user_id string) ([]domain.Server, error) {
	// Get all persistent servers first
	persistentServers, err := r.Persistent_repo.GetAllServers(ctx, user_id)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Server, 0, len(persistentServers))

	for _, ps := range persistentServers {
		// Try to get dynamic server data
		dynamicServers, err := r.Dynamic_repo.GetServerById(ctx, user_id, ps.ID)
		if err != nil {
			return nil, err
		}

		if len(dynamicServers) > 0 {
			// Dynamic record exists, use it
			ds := dynamicServers[0]
			result = append(result, domain.Server{
				ID:        ds.ID,
				UserId:    ds.UserId,
				Name:      ds.Name,
				Type:      ds.Type,
				Version:   ds.Version,
				CreatedAt: ds.CreatedAt,
				UpdatedAt: ds.UpdatedAt,
				DeletedAt: ds.DeletedAt,
				IP:        ds.IP,
				Status:    ds.Status,
				WorkerID:  ds.WorkerID,
			})
		} else {
			// No dynamic record, use persistent with stopped status
			result = append(result, domain.Server{
				ID:        ps.ID,
				UserId:    ps.UserId,
				Name:      ps.Name,
				Type:      ps.Type,
				Version:   ps.Version,
				CreatedAt: ps.CreatedAt,
				UpdatedAt: ps.UpdatedAt,
				DeletedAt: ps.DeletedAt,
				IP:        "",
				Status:    domain.StatusStopped,
				WorkerID:  "",
			})
		}
	}

	return result, nil
}

func (r *ServerRepository) GetServerById(ctx context.Context, user_id string, server_id string) ([]domain.Server, error) {
	// Try to get dynamic server data first
	dynamicServers, err := r.Dynamic_repo.GetServerById(ctx, user_id, server_id)
	if err != nil {
		return nil, err
	}

	if len(dynamicServers) > 0 {
		// Dynamic record exists, use it
		ds := dynamicServers[0]
		return []domain.Server{{
			ID:        ds.ID,
			UserId:    ds.UserId,
			Name:      ds.Name,
			Type:      ds.Type,
			Version:   ds.Version,
			CreatedAt: ds.CreatedAt,
			UpdatedAt: ds.UpdatedAt,
			DeletedAt: ds.DeletedAt,
			IP:        ds.IP,
			Status:    ds.Status,
			WorkerID:  ds.WorkerID,
		}}, nil
	}

	// No dynamic record, get persistent and return with stopped status
	persistentServers, err := r.Persistent_repo.GetServerById(ctx, user_id, server_id)
	if err != nil {
		return nil, err
	}

	if len(persistentServers) == 0 {
		return []domain.Server{}, nil
	}

	ps := persistentServers[0]
	return []domain.Server{{
		ID:        ps.ID,
		UserId:    ps.UserId,
		Name:      ps.Name,
		Type:      ps.Type,
		Version:   ps.Version,
		CreatedAt: ps.CreatedAt,
		UpdatedAt: ps.UpdatedAt,
		DeletedAt: ps.DeletedAt,
		IP:        "",
		Status:    domain.StatusStopped,
		WorkerID:  "",
	}}, nil
}

func (r *ServerRepository) CreateServer(ctx context.Context, server domain.ServerPersistent) error {
	return r.Persistent_repo.CreateServer(ctx, server)
}

func (r *ServerRepository) AddServerDynamic(ctx context.Context, user_id string, server_id string, ip string, status domain.ServerStatus, worker_id string) error {
	servers_persistent, err := r.Persistent_repo.GetServerById(ctx, user_id, server_id)
	if err != nil {
		return err
	}

	if len(servers_persistent) == 0 {
		return fmt.Errorf("server not found")
	}

	server_persistent := servers_persistent[0]

	server_dynamic := domain.ServerDynamic{
		ID:        server_persistent.ID,
		UserId:    user_id,
		Name:      server_persistent.Name,
		Type:      server_persistent.Type,
		Version:   server_persistent.Version,
		CreatedAt: server_persistent.CreatedAt,
		UpdatedAt: server_persistent.UpdatedAt,
		DeletedAt: server_persistent.DeletedAt,
		IP:        ip,
		Status:    status,
		WorkerID:  worker_id,
	}

	return r.Dynamic_repo.CreateServer(ctx, server_dynamic)
}
