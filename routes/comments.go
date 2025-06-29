package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"chirp/models"
)

// @Summary Создать комментарий
// @Description Создаёт новый комментарий к посту
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body routes.CreateCommentRequest true "Данные для комментария"
// @Success 201 {object} routes.CommentDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /comments [post]
func createCommentHandler(c *gin.Context, db *gorm.DB) {
	var req CreateCommentRequest
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

	comment := models.Comment{
		PostID:     req.PostID,
		AuthorID:   authorID,
		Content:    req.Content,
		IsReply:    req.ReplyToID != nil,
		ReplyToID:  req.ReplyToID,
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	resp := CommentDTO{
		ID:         comment.ID,
		PostID:     comment.PostID,
		AuthorID:   comment.AuthorID,
		Content:    comment.Content,
		Reputation: comment.Reputation,
		IsReply:    comment.IsReply,
		ReplyToID:  comment.ReplyToID,
		CreatedAt:  comment.CreatedAt,
	}

	c.JSON(http.StatusCreated, resp)
}

// @Summary Получить комментарии к посту
// @Description Получает все комментарии к посту
// @Tags comments
// @Produce json
// @Param id path string true "ID поста"
// @Success 200 {array} models.Comment
// @Failure 500 {object} map[string]string
// @Router /comments/posts/{id}/comments [get]
func getCommentsForPostHandler(c *gin.Context, db *gorm.DB) {
	postId := c.Param("id")
	var comments []models.Comment
	if err := db.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// @Summary Обновить комментарий
// @Description Обновляет комментарий пользователя
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID комментария"
// @Param data body routes.UpdateCommentDTO true "Данные для обновления"
// @Success 200 {object} models.Comment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /comments/{id} [put]
func updateCommentHandler(c *gin.Context, db *gorm.DB) {
	commentId := c.Param("id")
	var req UpdateCommentDTO
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

	var comment models.Comment
	if err := db.First(&comment,"id = ?", commentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to edit this comment"})
		return
	}

	comment.Content = req.Content
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary Удалить комментарий
// @Description Удаляет комментарий пользователя
// @Tags comments
// @Security BearerAuth
// @Param id path string true "ID комментария"
// @Success 204 {string} string ""
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /comments/{id} [delete]
func deleteCommentHandler(c *gin.Context, db *gorm.DB) {
	commentId := c.Param("id")

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

	var comment models.Comment
	if err := db.First(&comment, "id = ?", commentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this comment"})
		return
	}

	if err := db.Delete(&models.Comment{}, commentId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Голосовать за комментарий
// @Description Голосует за комментарий (лайк/дизлайк)
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID комментария"
// @Param data body routes.VoteDTO true "Голос"
// @Success 200 {object} models.Comment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /comments/{id}/vote [post]
func voteCommentHandler(c *gin.Context, db *gorm.DB) {
	commentId := c.Param("id")
	var req VoteDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var comment models.Comment
	if err := db.First(&comment, "id = ?" ,commentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	comment.Reputation += req.Value
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment reputation"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func RegisterCommentRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/", func(c *gin.Context) {
		createCommentHandler(c, db)
	})

	r.GET("/posts/:id/comments", func(c *gin.Context) {
		getCommentsForPostHandler(c, db)
	})

	r.PUT("/:id", func(c *gin.Context) {
		updateCommentHandler(c, db)
	})

	r.DELETE("/:id", func(c *gin.Context) {
		deleteCommentHandler(c, db)
	})

	r.POST("/:id/vote", func(c *gin.Context) {
		voteCommentHandler(c, db)
	})
}