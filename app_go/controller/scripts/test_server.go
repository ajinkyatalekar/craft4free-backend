/*
 * Prints out worker stats in the redis database
 */
 package main

 import (
	 "context"
	 "log"
 
	 "host-lotus-controller/internal/config"
	 "host-lotus-controller/internal/repository"
	 "host-lotus-controller/internal/service"
	 "host-lotus-controller/internal/domain"
 
	 "github.com/joho/godotenv"
 )
 
 func TestServer() {
 
	 _ = godotenv.Load()
 
	 // Load config
	 cfg, err := config.Load()
	 if err != nil {
		 log.Fatalf("Failed to load config: %v", err)
	 }
 
	 // Initialize services
	 redisClient, err := service.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	 if err != nil {
		 log.Fatalf("Failed to initialize Redis client: %v", err)
	 }
 
	 supabaseClient, err := service.NewSupabaseClient(cfg.SupabaseURL, cfg.SupabaseServiceKey)
	 if err != nil {
		 log.Fatalf("Failed to initialize Supabase client: %v", err)
	 }
 
	 serverRepo := repository.NewServerRepository(redisClient, supabaseClient)
 
	 ctx := context.Background()

	 err = serverRepo.AddServerDynamic(ctx, "51ef02ae-6ddc-437e-91ad-a2b292b10db4", "b081413f-b566-4192-b3e7-ac9879846f4d", "192.168.1.1", domain.StatusScheduled, "worker-uuid-or-identifier-2")
	if err != nil {
		log.Fatalf("Failed to add server dynamic: %v", err)
	}
}