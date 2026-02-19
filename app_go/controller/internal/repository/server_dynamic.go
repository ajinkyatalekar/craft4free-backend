package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"host-lotus-controller/internal/domain"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/supabase-community/supabase-go"
)

const (
	serverKey     = "server:dynamic:%s"
	allServersKey = "server:dynamic"
)

type serverDynamicRepository struct {
	supabase_client *supabase.Client
	redis_client    *redis.Client
}

func newServerDynamicRepository(supabase_client *supabase.Client, redis_client *redis.Client) *serverDynamicRepository {
	return &serverDynamicRepository{supabase_client: supabase_client, redis_client: redis_client}
}

func unmarshalServer(serverJson string) (domain.ServerDynamic, error) {
	var server domain.ServerDynamic
	err := json.Unmarshal([]byte(serverJson), &server)
	if err != nil {
		return domain.ServerDynamic{}, fmt.Errorf("failed to unmarshal server data: %w", err)
	}

	return server, nil
}

func (r *serverDynamicRepository) CreateServer(ctx context.Context, server domain.ServerDynamic) error {
	data, err := json.Marshal(server)
	if err != nil {
		return fmt.Errorf("failed to marshal server data: %w", err)
	}
	err = r.redis_client.Set(ctx, fmt.Sprintf(serverKey, server.ID), data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set server data: %w", err)
	}
	return err
}

func (r *serverDynamicRepository) GetServerById(ctx context.Context, user_id string, server_id string) ([]domain.ServerDynamic, error) {
	data, err := r.redis_client.Get(ctx, fmt.Sprintf(serverKey, server_id)).Result()
	if err != nil {
        if errors.Is(err, redis.Nil) {
            // Key doesn't exist - return empty slice or a custom "not found" error
            return []domain.ServerDynamic{}, nil
        }
        return nil, fmt.Errorf("failed to get server data: %w", err)
    }

	server, err := unmarshalServer(data)
	if err != nil {
		return nil, err
	}

	return []domain.ServerDynamic{server}, nil
}
