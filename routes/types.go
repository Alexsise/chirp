package routes

import (
	"time"

	"github.com/google/uuid"
)

// auth.go
// Представляет тело запроса для регистрации пользователя.
type RegisterRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Представляет ответ для зарегистрированного пользователя.
type RegisterResponse struct {
	ID           uuid.UUID `json:"id"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	RegisteredAt time.Time `json:"registeredAt"`
}

// Представляет тело запроса для входа в систему.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Представляет ответ с токеном после входа в систему.
type LoginResponse struct {
	Token string `json:"token"`
}

// Представляет ответ с обновлённым токеном.
type RefreshResponse struct {
	Token string `json:"token"`
}

// users.go
// Представляет приватный профиль пользователя.
type UserProfile struct {
	ID                uuid.UUID `json:"id"`
	Nickname          string    `json:"nickname"`
	Email             string    `json:"email"`
	BannerURL         string    `json:"bannerUrl"`
	PostReputation    int       `json:"postReputation"`
	CommentReputation int       `json:"commentReputation"`
	RegisteredAt      time.Time `json:"registeredAt"`
}

// Представляет тело запроса для обновления профиля пользователя.
type UpdateUserProfileRequest struct {
	Nickname  *string `json:"nickname"`
	BannerURL *string `json:"bannerUrl"`
	Password  *string `json:"password"`
}

// Представляет публичный профиль пользователя.
type PublicUserProfile struct {
	ID        uuid.UUID `json:"id"`
	Nickname  string    `json:"nickname"`
	BannerURL string    `json:"bannerUrl"`
}

// comments.go
// Представляет тело запроса для создания комментария.
type CreateCommentRequest struct {
	PostID    uuid.UUID `json:"postId" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	ReplyToID *uuid.UUID `json:"replyToId"`
}

// Представляет DTO для комментария.
type CommentDTO struct {
	ID         uuid.UUID `json:"id"`
	PostID     uuid.UUID `json:"postId"`
	AuthorID   uuid.UUID `json:"authorId"`
	Content    string    `json:"content"`
	Reputation int       `json:"reputation"`
	IsReply    bool      `json:"isReply"`
	ReplyToID  *uuid.UUID `json:"replyToId"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Представляет тело запроса для голосования за комментарий.
type VoteCommentRequest struct {
	Value int `json:"value" binding:"required,oneof=-1 1"`
}

// Представляет ответ на голосование за комментарий.
type VoteCommentResponse struct {
	Reputation int `json:"reputation"`
}

// Представляет тело запроса для обновления комментария.
type UpdateCommentDTO struct {
	Content string `json:"content" binding:"required"`
}

// Представляет тело запроса для голосования за комментарий.
type VoteDTO struct {
	Value int `json:"value" binding:"required,oneof=-1 1"`
}

// posts.go
// Представляет тело запроса для создания поста.
type CreatePostRequest struct {
	Content   string    `json:"content" binding:"required"`
	MediaUrls []string  `json:"mediaUrls"`
	GroupID   *uuid.UUID `json:"groupId"`
}

// Представляет DTO для поста.
type PostDTO struct {
	ID         uuid.UUID `json:"id"`
	AuthorID   uuid.UUID `json:"authorId"`
	Content    string    `json:"content"`
	MediaUrls  []string  `json:"mediaUrls"`
	Reputation int       `json:"reputation"`
	CreatedAt  time.Time `json:"createdAt"`
	GroupID    *uuid.UUID `json:"groupId"`
}

// Представляет ответ с постами с пагинацией.
type PaginatedPostsResponse struct {
	Posts      []PostDTO `json:"posts"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalCount int64     `json:"totalCount"`
}

// Представляет детализированный DTO для поста.
type PostDetailDTO struct {
	PostDTO
	Comments []CommentDTO `json:"comments"`
}

// Представляет тело запроса для обновления поста.
type UpdatePostRequest struct {
	Content   *string   `json:"content"`
	MediaUrls *[]string `json:"mediaUrls"`
}

// Представляет тело запроса для голосования за пост.
type VoteRequest struct {
	Value int `json:"value" binding:"required,oneof=-1 1"`
}

// Представляет ответ на голосование за пост.
type VoteResponse struct {
	Reputation int `json:"reputation"`
}

// groups.go
// Представляет тело запроса для создания группы.
type CreateGroupDTO struct {
	GroupName   string `json:"groupName" binding:"required"`
	Description string `json:"description"`
	BannerURL   string `json:"bannerUrl"`
}

// Представляет DTO для группы.
type GroupDTO struct {
	ID           uuid.UUID `json:"id"`
	GroupName    string    `json:"groupName"`
	RegisteredAt time.Time `json:"registeredAt"`
	BannerURL    string    `json:"bannerUrl"`
	Description  string    `json:"description"`
}

// Представляет детализированный DTO для группы.
type GroupDetailDTO struct {
	GroupDTO
	Moderators []UserProfile `json:"moderators"`
	Users      []UserProfile `json:"users"`
}

// Представляет тело запроса для обновления группы.
type UpdateGroupDTO struct {
	Description *string `json:"description"`
	BannerURL   *string `json:"bannerUrl"`
}

// subscriptions.go
// Представляет тело запроса для подписки на пользователя.
type SubscribeDTO struct {
	TargetUserID uuid.UUID `json:"targetUserId" binding:"required"`
}

// moderation.go
// Представляет тело запроса для добавления модератора в группу.
type AddModDTO struct {
	UserID uuid.UUID `json:"userId" binding:"required"`
}