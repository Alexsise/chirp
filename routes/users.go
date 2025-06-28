package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"chirp/models"

	"github.com/google/uuid"
)

func getUserProfileHandler(c *gin.Context, db *gorm.DB) {
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

	var user models.User
	if err := db.First(&user, "id = ?", authorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	profile := UserProfile{
		ID:                user.ID,
		Nickname:          user.Nickname,
		Email:             user.Email,
		BannerURL:         user.BannerURL,
		PostReputation:    user.ReputationPosts,
		CommentReputation: user.ReputationComments,
		RegisteredAt:      user.RegisteredAt,
	}

	c.JSON(http.StatusOK, profile)
}

func updateUserProfileHandler(c *gin.Context, db *gorm.DB) {
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

	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.First(&user, authorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.BannerURL != nil {
		user.BannerURL = *req.BannerURL
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	profile := UserProfile{
		ID:                user.ID,
		Nickname:          user.Nickname,
		Email:             user.Email,
		BannerURL:         user.BannerURL,
		PostReputation:    user.ReputationPosts,
		CommentReputation: user.ReputationComments,
		RegisteredAt:      user.RegisteredAt,
	}

	c.JSON(http.StatusOK, profile)
}

func getPublicUserProfileHandler(c *gin.Context, db *gorm.DB) {
	userID := c.Param("id")

	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	profile := PublicUserProfile{
		ID:        user.ID,
		Nickname:  user.Nickname,
		BannerURL: user.BannerURL,
	}

	c.JSON(http.StatusOK, profile)
}

func RegisterUserRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.GET("/me", func(c *gin.Context) {
		getUserProfileHandler(c, db)
	})

	r.PUT("/me", func(c *gin.Context) {
		updateUserProfileHandler(c, db)
	})

	r.GET("/:id", func(c *gin.Context) {
		getPublicUserProfileHandler(c, db)
	})
}