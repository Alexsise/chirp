package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"chirp/models"
)

func registerHandler(c *gin.Context, db *gorm.DB) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Nickname:     req.Nickname,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		RegisteredAt: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	resp := RegisterResponse{
		ID:           user.ID,
		Nickname:     user.Nickname,
		Email:        user.Email,
		RegisteredAt: user.RegisteredAt,
	}

	c.JSON(http.StatusCreated, resp)
}

func loginHandler(c *gin.Context, db *gorm.DB) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	signingKey := os.Getenv("JWT_SECRET")
	if signingKey == "" {
		signingKey = "default_secret" // Fallback for development
	}

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

func refreshHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	tokenString = tokenString[len("Bearer "):] // Remove "Bearer " prefix

	signingKey := os.Getenv("JWT_SECRET")
	if signingKey == "" {
		signingKey = "default_secret" // Fallback for development
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(signingKey), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userId"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": claims["userId"],
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err = newToken.SignedString([]byte(signingKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{Token: tokenString})
}

func RegisterAuthRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		registerHandler(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		loginHandler(c, db)
	})

	r.POST("/refresh", refreshHandler)
}