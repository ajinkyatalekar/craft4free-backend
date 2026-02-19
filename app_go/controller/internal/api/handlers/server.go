package handlers

import (
	"host-lotus-controller/internal/repository"
	"host-lotus-controller/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"host-lotus-controller/internal/service"
)

type ServerHandler struct {
	serverRepo *repository.ServerRepository
	workerRepo *repository.WorkerRepository
	serverService *service.ServerService
}

func NewServerHandler(serverRepo *repository.ServerRepository, workerRepo *repository.WorkerRepository, serverService *service.ServerService) *ServerHandler {
	return &ServerHandler{
		serverRepo: serverRepo,
		workerRepo: workerRepo,
		serverService: serverService,
	}
}

// CreateServer creates a new server in the persistent database.
// This does not start or send any RPCs to workers.
// TODO: Input validation
func (h *ServerHandler) CreateServer(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	
	s_name := c.PostForm("name")
	s_type := c.PostForm("type")
	s_version := c.PostForm("version")

	server := domain.ServerPersistent{
		ID: uuid.New().String(),
		UserId: user_id.(string),
		Name: s_name,
		Type: s_type,
		Version: s_version,
	}

	err := h.serverRepo.CreateServer(c.Request.Context(), server)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true, "server": server})
}

// GetAllUserServers returns all servers for a user.
// If dynamic row exists, sends dynamic data. Otherwise, sends persistent data.
func (h *ServerHandler) GetAllUserServers(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	servers, err := h.serverRepo.GetAllServers(c.Request.Context(), user_id.(string))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "servers": servers})
}

// GetServerById returns a server by ID.
// If dynamic row exists, sends dynamic data. Otherwise, sends persistent data.
func (h *ServerHandler) GetServerById(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	server_id := c.Param("server_id")

	server, err := h.serverRepo.GetServerById(c.Request.Context(), user_id.(string), server_id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "server": server})
}

// StartServer starts a server by ID.
func (h *ServerHandler) StartServer(c *gin.Context) {
    userID, _ := c.Get("user_id")
    serverID := c.Param("server_id")
    err := h.serverService.StartServer(c.Request.Context(), userID.(string), serverID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"success": true})
}

func (h *ServerHandler) StopServer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	serverID := c.Param("server_id")
	err := h.serverService.StopServer(c.Request.Context(), userID.(string), serverID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}