package services

import (
	"testing"

	"ledgerly/db"
	"ledgerly/models"

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
}

func TestPettyCashService_InsufficientFunds(t *testing.T) {
	setupTestDB()
	s := &PettyCashService{}

	// No seed funds - debit should fail
	tx := &models.PettyCashTransaction{
		Type:        models.TransactionTypeDebit,
		Amount:      100,
		Description: "Test Debit",
		UserID:      "test-user",
	}

	err := s.CreateTransaction(tx)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
}

func TestPettyCashService_SufficientFunds(t *testing.T) {
	setupTestDB()
	s := &PettyCashService{}

	// Seed funds first
	credit := &models.PettyCashTransaction{
		Type:        models.TransactionTypeCredit,
		Amount:      500,
		Description: "Seed",
		UserID:      "system",
	}
	err := s.CreateTransaction(credit)
	assert.NoError(t, err)

	// Now debit should succeed
	debit := &models.PettyCashTransaction{
		Type:        models.TransactionTypeDebit,
		Amount:      100,
		Description: "Test Debit",
		UserID:      "test-user",
	}
	err = s.CreateTransaction(debit)
	assert.NoError(t, err)

	// Verify balance
	balance, err := s.GetBalance()
	assert.NoError(t, err)
	assert.Equal(t, 400.0, balance)
}

func TestExpenseService_ZeroAmount(t *testing.T) {
	setupTestDB()
	s := &ExpenseService{}

	expense := &models.Expense{
		Title:    "Test",
		Amount:   0,
		Category: "Food",
		UserID:   "test-user",
	}

	err := s.CreateExpense(expense)
	assert.Error(t, err)
	assert.Equal(t, "amount must be greater than zero", err.Error())
}

func TestExpenseService_NegativeAmount(t *testing.T) {
	setupTestDB()
	s := &ExpenseService{}

	expense := &models.Expense{
		Title:    "Test",
		Amount:   -50,
		Category: "Food",
		UserID:   "test-user",
	}

	err := s.CreateExpense(expense)
	assert.Error(t, err)
	assert.Equal(t, "amount must be greater than zero", err.Error())
}

func TestExpenseService_EmptyCategory(t *testing.T) {
	setupTestDB()
	s := &ExpenseService{}

	expense := &models.Expense{
		Title:    "Test",
		Amount:   100,
		Category: "",
		UserID:   "test-user",
	}

	err := s.CreateExpense(expense)
	assert.Error(t, err)
	assert.Equal(t, "category is mandatory", err.Error())
}

func TestExpenseService_ValidExpense(t *testing.T) {
	setupTestDB()
	s := &ExpenseService{}

	expense := &models.Expense{
		Title:    "Lunch",
		Amount:   25.50,
		Category: "Food",
		UserID:   "test-user",
	}

	err := s.CreateExpense(expense)
	assert.NoError(t, err)
	assert.NotZero(t, expense.ID)
}
