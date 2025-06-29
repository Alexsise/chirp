package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"chirp/models"
)

func addModeratorHandler(c *gin.Context, db *gorm.DB) {
	groupID := c.Param("groupId")
	var req AddModDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	authorID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var group models.Group
	if err := db.Preload("Moderators").First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	isModerator := false
	for _, mod := range group.Moderators {
		if mod.ID == authorID {
			isModerator = true
			break
		}
	}

	if !isModerator {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to add moderators to this group"})
		return
	}

	if err := db.Model(&group).Association("Moderators").Append(&models.User{ID: req.UserID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add moderator"})
		return
	}

	c.Status(http.StatusNoContent)
}

func removeModeratorHandler(c *gin.Context, db *gorm.DB) {
	groupID := c.Param("groupId")
	userIDParam := c.Param("userId")

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	authorID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var group models.Group
	if err := db.Preload("Moderators").First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	isModerator := false
	for _, mod := range group.Moderators {
		if mod.ID == authorID {
			isModerator = true
			break
		}
	}

	if !isModerator {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to remove moderators from this group"})
		return
	}

	if err := db.Model(&group).Association("Moderators").Delete(&models.User{ID: targetUserID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove moderator"})
		return
	}

	c.Status(http.StatusNoContent)
}

func RegisterModerationRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/groups/:groupId/moderators", JWTMiddleware(), func(c *gin.Context) {
		addModeratorHandler(c, db)
	})

	r.DELETE("/groups/:groupId/moderators/:userId", JWTMiddleware(), func(c *gin.Context) {
		removeModeratorHandler(c, db)
	})
}