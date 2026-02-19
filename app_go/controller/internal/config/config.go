package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	API_Port int

	Dev bool

	SupabaseURL string
	SupabaseJWTSecret string
	SupabaseServiceKey string
	SupabaseAnonKey string

	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() (*Config, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	if supabaseURL == "" {
		return nil, errors.New("SUPABASE_URL environment variable is required")
	}

	supabaseJWTSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if supabaseJWTSecret == "" {
		return nil, errors.New("SUPABASE_JWT_SECRET environment variable is required")
	}

	supabaseServiceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	if supabaseServiceKey == "" {
		return nil, errors.New("SUPABASE_KEY environment variable is required")
	}

	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")
	if supabaseAnonKey == "" {
		return nil, errors.New("SUPABASE_ANON_KEY environment variable is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, errors.New("REDIS_ADDR environment variable is required")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		return nil, errors.New("REDIS_PASSWORD environment variable is required")
	}

	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		var err error
		redisDB, err = strconv.Atoi(dbStr)
		if err != nil {
			return nil, errors.New("REDIS_DB must be a valid integer")
		}
	}

	return &Config{
		Dev:               true,
		API_Port:          8000,
		SupabaseURL:       supabaseURL,
		SupabaseJWTSecret: supabaseJWTSecret,
		SupabaseServiceKey:       supabaseServiceKey,
		SupabaseAnonKey:       supabaseAnonKey,
		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
	}, nil
}