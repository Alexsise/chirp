package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"chirp/models"

	"github.com/google/uuid"
)

// @Summary Создать пост
// @Description Создаёт новый пост
// @Tags posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body routes.CreatePostRequest true "Данные для поста"
// @Success 201 {object} routes.PostDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /posts [post]
func createPostHandler(c *gin.Context, db *gorm.DB) {
	var req CreatePostRequest
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

	post := models.Post{
		AuthorID:  authorID,
		Content:   req.Content,
		MediaUrls:  req.MediaUrls, // Assuming single media URL for simplicity
		CreatedAt: time.Now(),
		GroupID:   req.GroupID,
	}

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	resp := PostDTO{
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		Content:   post.Content,
		MediaUrls: post.MediaUrls,
		Reputation: post.Reputation,
		CreatedAt: post.CreatedAt,
		GroupID:   post.GroupID,
	}

	c.JSON(http.StatusCreated, resp)
}

// @Summary Получить список постов
// @Description Получает посты с пагинацией
// @Tags posts
// @Produce json
// @Param page query int false "Страница"
// @Param limit query int false "Лимит"
// @Param sort query string false "Сортировка (createdAt|reputation)"
// @Success 200 {object} routes.PaginatedPostsResponse
// @Router /posts [get]
func getPaginatedPostsHandler(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "createdAt")
	if sort != "createdAt" && sort != "reputation" {
		sort = "createdAt" // Default to createdAt if invalid sort value is provided
	}

	offset := (page - 1) * limit

	var posts []models.Post
	var totalCount int64

	db.Model(&models.Post{}).Count(&totalCount)
	db.Order(sort + " DESC").Offset(offset).Limit(limit).Find(&posts)

	postDTOs := make([]PostDTO, len(posts))
	for i, post := range posts {
		postDTOs[i] = PostDTO{
			ID:        post.ID,
			AuthorID:  post.AuthorID,
			Content:   post.Content,
			MediaUrls: post.MediaUrls,
			Reputation: post.Reputation,
			CreatedAt: post.CreatedAt,
			GroupID:   post.GroupID,
		}
	}

	resp := PaginatedPostsResponse{
		Posts:      postDTOs,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Получить детали поста
// @Description Получает подробную информацию о посте
// @Tags posts
// @Produce json
// @Param id path string true "ID поста"
// @Success 200 {object} routes.PostDetailDTO
// @Failure 404 {object} map[string]string
// @Router /posts/{id} [get]
func getPostDetailHandler(c *gin.Context, db *gorm.DB) {
	postId := c.Param("postId")
	var post models.Post
	if err := db.Preload("Comments").First(&post, "id = ?",postId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comments := make([]CommentDTO, len(post.Comments))
	for i, comment := range post.Comments {
		comments[i] = CommentDTO{
			ID:         comment.ID,
			PostID:     comment.PostID,
			AuthorID:   comment.AuthorID,
			Content:    comment.Content,
			Reputation: comment.Reputation,
			IsReply:    comment.IsReply,
			ReplyToID:  comment.ReplyToID,
			CreatedAt:  time.Now(),
		}
	}

	resp := PostDetailDTO{
		PostDTO: PostDTO{
			ID:        post.ID,
			AuthorID:  post.AuthorID,
			Content:   post.Content,
			MediaUrls: post.MediaUrls,
			Reputation: post.Reputation,
			CreatedAt: time.Now(),
			GroupID:   post.GroupID,
		},
		Comments: comments,
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Обновить пост
// @Description Обновляет пост пользователя
// @Tags posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID поста"
// @Param data body routes.UpdatePostRequest true "Данные для обновления"
// @Success 200 {object} routes.PostDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /posts/{id} [put]
func updatePostHandler(c *gin.Context, db *gorm.DB) {
	postId := c.Param("postId")
	var req UpdatePostRequest
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

	var post models.Post
	if err := db.First(&post, "id = ?", postId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to edit this post"})
		return
	}

	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.MediaUrls != nil && len(*req.MediaUrls) > 0 {
		post.MediaUrls = *req.MediaUrls
	}

	if err := db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	resp := PostDTO{
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		Content:   post.Content,
		MediaUrls: post.MediaUrls,
		Reputation: post.Reputation,
		CreatedAt: post.CreatedAt,
		GroupID:   post.GroupID,
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Удалить пост
// @Description Удаляет пост пользователя
// @Tags posts
// @Security BearerAuth
// @Param id path string true "ID поста"
// @Success 204 {string} string ""
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /posts/{id} [delete]
func deletePostHandler(c *gin.Context, db *gorm.DB) {
	postId := c.Param("postId")

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

	var post models.Post
	if err := db.First(&post, "id = ?", postId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.AuthorID != authorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this post"})
		return
	}

	if err := db.Delete(&models.Post{}, postId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Голосовать за пост
// @Description Голосует за пост (лайк/дизлайк)
// @Tags posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID поста"
// @Param data body routes.VoteRequest true "Голос"
// @Success 200 {object} routes.VoteResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /posts/{id}/vote [post]
func votePostHandler(c *gin.Context, db *gorm.DB) {
	postId := c.Param("postId")
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post models.Post
	if err := db.First(&post, "id = ?",postId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	post.Reputation += req.Value
	if err := db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post reputation"})
		return
	}

	resp := VoteResponse{
		Reputation: post.Reputation,
	}

	c.JSON(http.StatusOK, resp)
}

func RegisterPostRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/", JWTMiddleware(), func(c *gin.Context) {
		createPostHandler(c, db)
	})

	r.GET("/", func(c *gin.Context) {
		getPaginatedPostsHandler(c, db)
	})

	r.GET("/:id", func(c *gin.Context) {
		getPostDetailHandler(c, db)
	})

	r.PUT("/:id", JWTMiddleware(), func(c *gin.Context) {
		updatePostHandler(c, db)
	})

	r.DELETE("/:id", JWTMiddleware(), func(c *gin.Context) {
		deletePostHandler(c, db)
	})

	r.POST("/:id/vote", JWTMiddleware(), func(c *gin.Context) {
		votePostHandler(c, db)
	})
}