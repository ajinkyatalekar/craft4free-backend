package routes

import (
	"github.com/gin-gonic/gin"

	"host-lotus-controller/internal/api/middleware"
	"host-lotus-controller/internal/api/handlers"
	"host-lotus-controller/internal/config"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, serverHandler *handlers.ServerHandler) {
	// Public routes
	r.GET("/health", handlers.PublicHealth)

	// Protected routes (auth required)
	protected := r.Group("/api")
	protected.Use(middleware.SupabaseAuth(cfg.SupabaseJWTSecret))
	{
		protected.GET("/health", func(c *gin.Context) { handlers.ProtectedHealth(c) })
		protected.GET("/servers", serverHandler.GetAllUserServers)
		protected.GET("/servers/:server_id", serverHandler.GetServerById)
		protected.POST("/servers", serverHandler.CreateServer)
		protected.POST("/servers/:server_id/start", serverHandler.StartServer)
		protected.POST("/servers/:server_id/stop", serverHandler.StopServer)
	}
}