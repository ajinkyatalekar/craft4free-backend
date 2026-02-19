/*
 * Prints out worker stats in the redis database
 */
package main

import (
	"context"
	"fmt"
	"log"

	"host-lotus-controller/internal/config"
	"host-lotus-controller/internal/repository"
	"host-lotus-controller/internal/service"

	"github.com/joho/godotenv"
)

func main() {

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

	workerRepo := repository.NewWorkerRepository(redisClient, supabaseClient)

	ctx := context.Background()

	// err = workerRepo.AssignServerToWorker(ctx, "worker-uuid-or-identifier", "t-1")
	// if err != nil {
	// 	log.Printf("Error assigning server: %v", err)
	// } else {
	// 	fmt.Println("Successfully assigned server to worker")
	// }

	// err = workerRepo.AssignServer(ctx, "t-1")
	// if err != nil {
	// 	log.Printf("Error assigning server: %v", err)
	// } else {
	// 	fmt.Println("Successfully assigned server to worker")
	// }

	workers, err := workerRepo.GetAllWorkers(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Workers found:\n")
		for _, worker := range workers {
			fmt.Printf("  ID: %s\n", worker.ID)
		}
	}

	onlineWorkers, err := workerRepo.GetOnlineWorkers(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Online workers found:\n")
		for _, worker := range onlineWorkers {
			fmt.Printf("  ID: %s\n", worker.ID)
		}
	}

	freeWorkers, err := workerRepo.GetFreeWorkers(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Free workers found:\n")
		for _, worker := range freeWorkers {
			fmt.Printf("  ID: %s\n", worker.ID)
			fmt.Printf("  Free slots: %d\n", worker.MaxServers-len(worker.RunningServers)-len(worker.AssignedServers))
		}
	}
}
