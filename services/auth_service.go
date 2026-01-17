package services

import (
	"errors"
	"log/slog"
	"time"
	"ledgerly/db"
	"ledgerly/middleware"
	"ledgerly/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

var jwtKey = []byte(middleware.GetJWTSecret())

func (s *AuthService) Register(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return err
	}
	user.Password = string(hashedPassword)
	return db.DB.Create(user).Error
}

func (s *AuthService) Login(username, password string) (string, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		slog.Warn("Login failed: invalid username", "username", username)
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		slog.Warn("Login failed: invalid password", "username", username)
		return "", errors.New("invalid credentials")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		slog.Error("Failed to sign token", "error", err)
		return "", err
	}

	slog.Info("User logged in successfully", "username", username, "role", user.Role)
	return tokenString, nil
}
