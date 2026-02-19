package repository

import (
	"context"
	"encoding/json"
	"host-lotus-controller/internal/domain"

	"github.com/redis/go-redis/v9"
	"github.com/supabase-community/supabase-go"
)

type serverPersistentRepository struct {
	supabase_client *supabase.Client
	redis_client    *redis.Client
}

func newServerPersistentRepository(supabase_client *supabase.Client, redis_client *redis.Client) *serverPersistentRepository {
	return &serverPersistentRepository{supabase_client: supabase_client, redis_client: redis_client}
}

func (r *serverPersistentRepository) GetAllServers(ctx context.Context, user_id string) ([]domain.ServerPersistent, error) {
	var servers []domain.ServerPersistent

	// Select all columns (*) from the servers table
	data, _, err := r.supabase_client.From("server").
		Select("*", "exact", false).
		Eq("user_id", user_id).
		Execute()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *serverPersistentRepository) GetServerById(ctx context.Context, user_id string, server_id string) ([]domain.ServerPersistent, error) {
	var servers []domain.ServerPersistent

	// Select all columns (*) from the servers table
	data, _, err := r.supabase_client.From("server").
		Select("*", "exact", false).
		Eq("user_id", user_id).
		Eq("id", server_id).
		Execute()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *serverPersistentRepository) CreateServer(ctx context.Context, server domain.ServerPersistent) error {
	// Insert the server into the servers table
	_, _, err := r.supabase_client.From("server").
		Insert(server, false, "", "representation", "exact").
		Execute()
	if err != nil {
		return err
	}

	return nil
}
