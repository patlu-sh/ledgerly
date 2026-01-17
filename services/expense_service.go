package services

import (
	"errors"
	"ledgerly/db"
	"ledgerly/models"
)

type ExpenseService struct{}

func (s *ExpenseService) CreateExpense(expense *models.Expense) error {
	if expense.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if expense.Category == "" {
		return errors.New("category is mandatory")
	}
	return db.DB.Create(expense).Error
}

func (s *ExpenseService) ListExpenses() ([]models.Expense, error) {
	var expenses []models.Expense
	err := db.DB.Preload("PettyCashTransaction").Find(&expenses).Error
	return expenses, err
}
