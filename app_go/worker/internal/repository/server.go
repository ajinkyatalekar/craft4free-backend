package repository

import (
	"context"	
	"encoding/json"
	"fmt"
	"host-lotus-worker/internal/domain"
	"github.com/redis/go-redis/v9"
	"github.com/supabase-community/supabase-go"
)

const (
	serverKey     = "server:dynamic:%s"
	allServersKey = "server:dynamic"
)

type ServerRepository struct {
	redis_client *redis.Client
	supabase_client *supabase.Client
}

func NewServerRepository(redis_client *redis.Client, supabase_client *supabase.Client) *ServerRepository {
	return &ServerRepository{
		redis_client: redis_client,
		supabase_client: supabase_client,
	}
}

func unmarshalServer(serverJson string) (domain.ServerDynamic, error) {
	var server domain.ServerDynamic
	err := json.Unmarshal([]byte(serverJson), &server)
	if err != nil {
		return domain.ServerDynamic{}, fmt.Errorf("failed to unmarshal server data: %w", err)
	}

	return server, nil
}

func (r *ServerRepository) GetServerById(ctx context.Context, server_id string) (domain.ServerDynamic, error) {
	data, err := r.redis_client.Get(ctx, fmt.Sprintf(serverKey, server_id)).Result()
	if err != nil {
		return domain.ServerDynamic{}, fmt.Errorf("failed to get server data: %w", err)
	}

	server, err := unmarshalServer(data)
	if err != nil {
		return domain.ServerDynamic{}, fmt.Errorf("failed to unmarshal server data: %w", err)
	}

	return server, nil
}

func (r *ServerRepository) UpdateServer(ctx context.Context, server domain.ServerDynamic) error {
	data, err := json.Marshal(server)
	if err != nil {
		return fmt.Errorf("failed to marshal server data: %w", err)
	}
	err = r.redis_client.Set(ctx, fmt.Sprintf(serverKey, server.ID), data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to update server data: %w", err)
	}
	return nil
}