package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"ledgerly/db"
	"ledgerly/middleware"
	"ledgerly/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	db.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.DB.AutoMigrate(&models.PettyCashTransaction{}, &models.Expense{}, &models.User{})
	
	// Seed initial funds for Debit tests
	db.DB.Create(&models.PettyCashTransaction{
		Type:        models.TransactionTypeCredit,
		Amount:      10000,
		Description: "Initial Seed",
		UserID:      "system",
	})
}

func generateToken(role models.UserRole) string {
	claims := &middleware.Claims{
		UserID: 1,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Use the default secret from auth_middleware.go if env is not set
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret_key"
	}
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestRBAC_Routes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	r := SetupRouter()

	adminToken := generateToken(models.RoleAdmin)
	employeeToken := generateToken(models.RoleEmployee)

	tests := []struct {
		name           string
		method         string
		url            string
		token          string
		body           interface{}
		expectedStatus int
	}{
		// Petty Cash Create (The Problematic Route)
		{
			name:   "Admin Create Petty Cash (Should succeed now)",
			method: "POST",
			url:    "/petty-cash",
			token:  adminToken,
			body: models.PettyCashTransaction{
				Type:   models.TransactionTypeCredit,
				Amount: 100,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Employee Create Petty Cash Debit (Should succeed)",
			method: "POST",
			url:    "/petty-cash",
			token:  employeeToken,
			body: models.PettyCashTransaction{
				Type:   models.TransactionTypeDebit,
				Amount: 50,
			},
			expectedStatus: http.StatusCreated, // Employee has expenses.create -> mapped to Debit
		},

		// Petty Cash List
		{
			name:           "Admin List Petty Cash",
			method:         "GET",
			url:            "/petty-cash",
			token:          adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Employee List Petty Cash",
			method:         "GET",
			url:            "/petty-cash",
			token:          employeeToken,
			expectedStatus: http.StatusForbidden, // Employee doesn't have view_list
		},

		// Petty Cash Balance
		{
			name:           "Admin View Balance",
			method:         "GET",
			url:            "/petty-cash/balance",
			token:          adminToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Employee View Balance",
			method:         "GET",
			url:            "/petty-cash/balance",
			token:          employeeToken,
			expectedStatus: http.StatusOK,
		},

		// Expenses Create
		{
			name:   "Admin Create Expense",
			method: "POST",
			url:    "/expenses",
			token:  adminToken,
			body: models.Expense{
				Title:    "Test",
				Amount:   100,
				Category: "Food",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Employee Create Expense",
			method: "POST",
			url:    "/expenses",
			token:  employeeToken,
			body: models.Expense{
				Title:    "Test",
				Amount:   100,
				Category: "Food",
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
