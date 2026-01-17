package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ledgerly/db"
	"ledgerly/models"
	"ledgerly/services"

	"github.com/gin-gonic/gin"
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
	
	// Seed initial funds
	db.DB.Create(&models.PettyCashTransaction{
		Type:        models.TransactionTypeCredit,
		Amount:      10000,
		Description: "Initial Seed",
		UserID:      "system",
	})
}

func TestCreatePettyCashTransaction_RBAC(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	tests := []struct {
		name            string
		role            models.UserRole
		transactionType models.TransactionType
		expectedStatus  int
	}{
		{
			name:            "Admin Create Credit",
			role:            models.RoleAdmin,
			transactionType: models.TransactionTypeCredit,
			expectedStatus:  http.StatusCreated,
		},
		{
			name:            "Admin Create Debit",
			role:            models.RoleAdmin,
			transactionType: models.TransactionTypeDebit,
			expectedStatus:  http.StatusCreated,
		},
		{
			name:            "Employee Create Credit",
			role:            models.RoleEmployee,
			transactionType: models.TransactionTypeCredit,
			expectedStatus:  http.StatusForbidden,
		},
		{
			name:            "Employee Create Debit",
			role:            models.RoleEmployee,
			transactionType: models.TransactionTypeDebit,
			expectedStatus:  http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tx := models.PettyCashTransaction{
				Type:        tt.transactionType,
				Amount:      100,
				Description: "Test Transaction",
				UserID:      "test-user",
			}
			jsonBytes, _ := json.Marshal(tx)
			c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			c.Set("role", tt.role)

			h := &Handler{
				PettyCashService: &services.PettyCashService{},
			}

			h.CreatePettyCashTransaction(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
