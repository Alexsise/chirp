package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"chirp/models"
)

func subscribeToGroupHandler(c *gin.Context, db *gorm.DB) {
	groupID := c.Param("groupId")

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	authorID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var group models.Group
	if err := db.Preload("Users").First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	for _, user := range group.Users {
		if user.ID == authorID {
			c.Status(http.StatusNoContent) // Already subscribed
			return
		}
	}

	if err := db.Model(&group).Association("Users").Append(&models.User{ID: authorID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe to group"})
		return
	}

	c.Status(http.StatusNoContent)
}

func unsubscribeFromGroupHandler(c *gin.Context, db *gorm.DB) {
	groupID := c.Param("groupId")

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	authorID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var group models.Group
	if err := db.Preload("Users").First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	if err := db.Model(&group).Association("Users").Delete(&models.User{ID: authorID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe from group"})
		return
	}

	c.Status(http.StatusNoContent)
}

func RegisterSubscriptionRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/groups/:groupId/subscribe", JWTMiddleware(), func(c *gin.Context) {
		subscribeToGroupHandler(c, db)
	})

	r.DELETE("/groups/:groupId/subscribe", JWTMiddleware(), func(c *gin.Context) {
		unsubscribeFromGroupHandler(c, db)
	})
}