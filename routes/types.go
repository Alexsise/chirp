package routes

import (
	"time"

	"github.com/google/uuid"
)

// auth.go
// Represents the request payload for registering a user.
type RegisterRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Represents the response payload for a registered user.
type RegisterResponse struct {
	ID           uuid.UUID `json:"id"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	RegisteredAt time.Time `json:"registeredAt"`
}

// Represents the request payload for logging in.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Represents the response payload for a login token.
type LoginResponse struct {
	Token string `json:"token"`
}

// Represents the response payload for a refreshed token.
type RefreshResponse struct {
	Token string `json:"token"`
}

// users.go
// Represents the private user profile.
type UserProfile struct {
	ID                uuid.UUID `json:"id"`
	Nickname          string    `json:"nickname"`
	Email             string    `json:"email"`
	BannerURL         string    `json:"bannerUrl"`
	PostReputation    int       `json:"postReputation"`
	CommentReputation int       `json:"commentReputation"`
	RegisteredAt      time.Time `json:"registeredAt"`
}

// Represents the request payload for updating a user profile.
type UpdateUserProfileRequest struct {
	Nickname  *string `json:"nickname"`
	BannerURL *string `json:"bannerUrl"`
	Password  *string `json:"password"`
}

// Represents the public user profile.
type PublicUserProfile struct {
	ID        uuid.UUID `json:"id"`
	Nickname  string    `json:"nickname"`
	BannerURL string    `json:"bannerUrl"`
}

// comments.go
// Represents the request payload for creating a comment.
type CreateCommentRequest struct {
	PostID    uuid.UUID `json:"postId" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	ReplyToID *uuid.UUID `json:"replyToId"`
}

// Represents the data transfer object for a comment.
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

// Represents the request payload for voting on a comment.
type VoteCommentRequest struct {
	Value int `json:"value" binding:"required,oneof=-1 1"`
}

// Represents the response payload for voting on a comment.
type VoteCommentResponse struct {
	Reputation int `json:"reputation"`
}

// Represents the request payload for updating a comment.
type UpdateCommentDTO struct {
	Content string `json:"content" binding:"required"`
}

// Represents the request payload for voting on a comment.
type VoteDTO struct {
	Value int `json:"value" binding:"required,oneof=-1 1"`
}

// posts.go
// Represents the request payload for creating a post.
type CreatePostRequest struct {
	Content   string    `json:"content" binding:"required"`
	MediaUrls []string  `json:"mediaUrls"`
	GroupID   *uuid.UUID `json:"groupId"`
}

// Represents the data transfer object for a post.
type PostDTO struct {
	ID         uuid.UUID `json:"id"`
	AuthorID   uuid.UUID `json:"authorId"`
	Content    string    `json:"content"`
	MediaUrls  []string  `json:"mediaUrls"`
	Reputation int       `json:"reputation"`
	CreatedAt  time.Time `json:"createdAt"`
	GroupID    *uuid.UUID `json:"groupId"`
}

// Represents the response payload for paginated posts.
type PaginatedPostsResponse struct {
	Posts      []PostDTO `json:"posts"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalCount int64     `json:"totalCount"`
}

// Represents the detailed data transfer object for a post.
type PostDetailDTO struct {
	PostDTO
	Comments []CommentDTO `json:"comments"`
}

// Represents the request payload for updating a post.
type UpdatePostRequest struct {
	Content   *string   `json:"content"`
	MediaUrls *[]string `json:"mediaUrls"`
}

// Represents the request payload for voting on a post.
type VoteRequest struct {
	Value int `json:"value" binding:"required,oneof=-1 1"`
}

// Represents the response payload for voting on a post.
type VoteResponse struct {
	Reputation int `json:"reputation"`
}

// groups.go
// Represents the request payload for creating a group.
type CreateGroupDTO struct {
	GroupName   string `json:"groupName" binding:"required"`
	Description string `json:"description"`
	BannerURL   string `json:"bannerUrl"`
}

// Represents the data transfer object for a group.
type GroupDTO struct {
	ID           uuid.UUID `json:"id"`
	GroupName    string    `json:"groupName"`
	RegisteredAt time.Time `json:"registeredAt"`
	BannerURL    string    `json:"bannerUrl"`
	Description  string    `json:"description"`
}

// Represents the detailed data transfer object for a group.
type GroupDetailDTO struct {
	GroupDTO
	Moderators []UserProfile `json:"moderators"`
	Users      []UserProfile `json:"users"`
}

// Represents the request payload for updating a group.
type UpdateGroupDTO struct {
	Description *string `json:"description"`
	BannerURL   *string `json:"bannerUrl"`
}

// subscriptions.go
// Represents the request payload for subscribing to a user.
type SubscribeDTO struct {
	TargetUserID uuid.UUID `json:"targetUserId" binding:"required"`
}

// moderation.go
// Represents the request payload for adding a moderator to a group.
type AddModDTO struct {
	UserID uuid.UUID `json:"userId" binding:"required"`
}