package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	authGroup := r.Group("/api/v1/auth")
	RegisterAuthRoutes(authGroup, db)

	usersGroup := r.Group("/api/v1/users")
	RegisterUserRoutes(usersGroup, db)

	postsGroup := r.Group("/api/v1/posts")
	RegisterPostRoutes(postsGroup, db)

	commentsGroup := r.Group("/api/v1/comments")
	RegisterCommentRoutes(commentsGroup, db)
}