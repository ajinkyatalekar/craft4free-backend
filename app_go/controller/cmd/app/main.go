package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"host-lotus-controller/internal/api/routes"
	"host-lotus-controller/internal/config"
	"host-lotus-controller/internal/repository"
	"host-lotus-controller/internal/api/handlers"
	"host-lotus-controller/internal/service"
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
	
	// Initialize Redis client
	redisClient, err := service.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	serverRepo := repository.NewServerRepository(redisClient, supabaseClient)
	workerRepo := repository.NewWorkerRepository(redisClient, supabaseClient)

	// Initialize services
	serverService := service.NewServerService(serverRepo, workerRepo)

	// Initialize handlers with dependencies
	serverHandler := handlers.NewServerHandler(serverRepo, workerRepo, serverService)


	// Setup router
	gin.SetMode(gin.DebugMode)
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	routes.SetupRoutes(r, cfg, serverHandler)
	r.Run(fmt.Sprintf(":%d", cfg.API_Port))
}
