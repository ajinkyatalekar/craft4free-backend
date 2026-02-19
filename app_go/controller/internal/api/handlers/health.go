package handlers

import (
	"github.com/gin-gonic/gin"
)

func PublicHealth(c *gin.Context) {
	c.JSON(200, gin.H{"success": true, "message": "Controller is running"})
}

func ProtectedHealth(c *gin.Context) {
	c.JSON(200, gin.H{"success": true, "message": "Controller is running"})
}