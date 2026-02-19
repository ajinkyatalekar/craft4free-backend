package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"host-lotus-worker/internal/config"
	"host-lotus-worker/internal/domain"
	"host-lotus-worker/internal/repository"
	"host-lotus-worker/internal/service"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Supabase client
	supabaseClient, err := service.NewSupabaseClient(cfg.SupabaseURL, cfg.SupabaseServiceKey)
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}
	log.Println("Supabase client initialized successfully")

	// Initialize Redis client
	redisClient, err := service.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()
	log.Println("Redis client initialized successfully")

	serverRepo := repository.NewServerRepository(redisClient, supabaseClient)

	worker := &domain.Worker{
		ID:              cfg.ID,
		OciId:           cfg.OciId,
		CreatedAt:       cfg.CreatedAt,
		PublicIp:        cfg.PublicIp,
		PrivateIp:       cfg.PrivateIp,
		MaxServers:      cfg.MaxServers,
		RunningServers:  map[string]string{},
		AssignedServers: map[string]string{},
	}

	// Initialize and start heartbeat service
	heartbeatInterval := 1 * time.Second
	heartbeatService := service.NewHeartbeatService(redisClient, worker, heartbeatInterval)
	heartbeatService.Start()

	assignmentInterval := 500 * time.Millisecond
	assignmentProcessor := service.NewAssignmentProcessor(redisClient, worker, serverRepo, assignmentInterval)
	assignmentProcessor.Start()

	// Add command processor
	commandProcessor := service.NewCommandProcessor(redisClient, worker, serverRepo)
	commandProcessor.Start()

	log.Printf("Worker %s started successfully", worker.ID)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down worker...")
	heartbeatService.Stop()
	assignmentProcessor.Stop() // Add this if AssignmentProcessor has Stop()
	commandProcessor.Stop()    // Stop the command processor
	_ = supabaseClient
	log.Println("Worker shut down successfully")
}
