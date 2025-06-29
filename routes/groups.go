package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"chirp/models"
)

// @Summary Создать группу
// @Description Создаёт новую группу
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body routes.CreateGroupDTO true "Данные для группы"
// @Success 201 {object} routes.GroupDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /groups [post]
func createGroupHandler(c *gin.Context, db *gorm.DB) {
	var req CreateGroupDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group data"})
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

	group := models.Group{
		ID:           uuid.New(),
		GroupName:    req.GroupName,
		Description:  req.Description,
		BannerURL:    req.BannerURL,
		RegisteredAt: time.Now(),
		Moderators:   []models.User{{ID: authorID}},
	}

	if err := db.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	resp := GroupDTO{
		ID:           group.ID,
		GroupName:    group.GroupName,
		RegisteredAt: group.RegisteredAt,
		BannerURL:    group.BannerURL,
		Description:  group.Description,
	}

	c.JSON(http.StatusCreated, resp)
}

// @Summary Получить список групп
// @Description Получает список всех групп
// @Tags groups
// @Produce json
// @Success 200 {array} routes.GroupDTO
// @Failure 500 {object} map[string]string
// @Router /groups [get]
func listGroupsHandler(c *gin.Context, db *gorm.DB) {
	var groups []models.Group
	if err := db.Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}

	groupDTOs := make([]GroupDTO, len(groups))
	for i, group := range groups {
		groupDTOs[i] = GroupDTO{
			ID:           group.ID,
			GroupName:    group.GroupName,
			RegisteredAt: group.RegisteredAt,
			BannerURL:    group.BannerURL,
			Description:  group.Description,
		}
	}

	c.JSON(http.StatusOK, groupDTOs)
}

// @Summary Получить детали группы
// @Description Получает подробную информацию о группе
// @Tags groups
// @Produce json
// @Param id path string true "ID группы"
// @Success 200 {object} routes.GroupDetailDTO
// @Failure 404 {object} map[string]string
// @Router /groups/{id} [get]
func getGroupDetailsHandler(c *gin.Context, db *gorm.DB) {
	groupID := c.Param("id")

	var group models.Group
	if err := db.Preload("Moderators").Preload("Users").First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to find group"})
		return
	}

	moderators := make([]UserProfile, len(group.Moderators))
	for i, mod := range group.Moderators {
		moderators[i] = UserProfile{
			ID:       mod.ID,
			Nickname: mod.Nickname,
			BannerURL: mod.BannerURL,
		}
	}

	users := make([]UserProfile, len(group.Users))
	for i, user := range group.Users {
		users[i] = UserProfile{
			ID:       user.ID,
			Nickname: user.Nickname,
			BannerURL: user.BannerURL,
		}
	}

	resp := GroupDetailDTO{
		GroupDTO: GroupDTO{
			ID:           group.ID,
			GroupName:    group.GroupName,
			RegisteredAt: group.RegisteredAt,
			BannerURL:    group.BannerURL,
			Description:  group.Description,
		},
		Moderators: moderators,
		Users:      users,
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Обновить группу
// @Description Обновляет данные группы
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID группы"
// @Param data body routes.UpdateGroupDTO true "Данные для обновления"
// @Success 200 {object} routes.GroupDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /groups/{id} [put]
func updateGroupHandler(c *gin.Context, db *gorm.DB) {
	groupID := c.Param("id")
	var req UpdateGroupDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group data"})
		return
	}

	var group models.Group
	if err := db.First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to find group"})
		return
	}

	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.BannerURL != nil {
		group.BannerURL = *req.BannerURL
	}

	if err := db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group"})
		return
	}

	resp := GroupDTO{
		ID:           group.ID,
		GroupName:    group.GroupName,
		RegisteredAt: group.RegisteredAt,
		BannerURL:    group.BannerURL,
		Description:  group.Description,
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Удалить группу
// @Description Удаляет группу
// @Tags groups
// @Security BearerAuth
// @Param id path string true "ID группы"
// @Success 204 {string} string ""
// @Failure 500 {object} map[string]string
// @Router /groups/{id} [delete]
func deleteGroupHandler(c *gin.Context, db *gorm.DB) {
	groupID := c.Param("id")

	if err := db.Delete(&models.Group{}, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	c.Status(http.StatusNoContent)
}

func RegisterGroupRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/", JWTMiddleware(), func(c *gin.Context) {
		createGroupHandler(c, db)
	})

	r.GET("/", func(c *gin.Context) {
		listGroupsHandler(c, db)
	})

	r.GET("/:id", func(c *gin.Context) {
		getGroupDetailsHandler(c, db)
	})

	r.PUT("/:id", JWTMiddleware(), func(c *gin.Context) {
		updateGroupHandler(c, db)
	})

	r.DELETE("/:id", JWTMiddleware(), func(c *gin.Context) {
		deleteGroupHandler(c, db)
	})
}