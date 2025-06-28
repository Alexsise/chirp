package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Nickname           string    `gorm:"not null;unique"`
	Email              string    `gorm:"not null;unique"`
	ReputationPosts    int       `gorm:"default:0"`
	ReputationComments int       `gorm:"default:0"`
	PasswordHash       string    `gorm:"not null"`
	RegisteredAt       time.Time `gorm:"not null"`
	BannerURL          string
	Groups             []Group   `gorm:"many2many:group_users"`
	Subscriptions       []User    `gorm:"many2many:user_subscriptions;joinForeignKey:subscriber_id;joinReferences:target_user_id"`
}

type Post struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AuthorID   uuid.UUID `gorm:"type:uuid;not null"`
	Author     User      `gorm:"foreignKey:AuthorID"`
	Content    string    `gorm:"type:text;not null"`
	MediaUrls  string    `gorm:"type:json"`
	Reputation int       `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"not null"`
	GroupID    *uuid.UUID
	Group      *Group
	Comments   []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PostID     uuid.UUID `gorm:"type:uuid;not null"`
	AuthorID   uuid.UUID `gorm:"type:uuid;not null"`
	Author     User      `gorm:"foreignKey:AuthorID"`
	Content    string    `gorm:"type:text;not null"`
	Reputation int       `gorm:"default:0"`
	IsReply    bool      `gorm:"not null"`
	ReplyToID  *uuid.UUID
	CreatedAt  time.Time `gorm:"not null"`
}

type Group struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GroupName    string    `gorm:"not null;unique"`
	RegisteredAt time.Time `gorm:"not null"`
	BannerURL    string
	Description  string    `gorm:"type:text"`
	Moderators   []User    `gorm:"many2many:group_moderators"`
	Users        []User    `gorm:"many2many:group_users"`
}

type GroupUser struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GroupID  uuid.UUID `gorm:"type:uuid;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	JoinedAt time.Time `gorm:"not null"`
	Title    string    `gorm:"type:varchar(255)"`
}

type GroupModerator struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GroupID uuid.UUID `gorm:"type:uuid;not null"`
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
}

func (GroupModerator) TableName() string {
	return "group_moderators"
}

func InitDB() *gorm.DB {
	dsn := 
		"host=localhost user=chirp_user password=chirp_password dbname=chirp_db port=5432 sslmode=disable"
	

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&User{},
		&Post{},
		&Comment{},
		&Group{},
		&GroupUser{},
		&GroupModerator{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	return db
}
