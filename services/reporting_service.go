package services

import (
	"ledgerly/db"
	"ledgerly/models"
)

type ReportingService struct{}

type ExpenseSummary struct {
	TotalExpenses float64            `json:"total_expenses"`
	ByCategory    map[string]float64 `json:"by_category"`
}

type PettyCashSummary struct {
	TotalCredits float64 `json:"total_credits"`
	TotalDebits  float64 `json:"total_debits"`
	Balance      float64 `json:"balance"`
}

func (s *ReportingService) GetExpenseSummary() (*ExpenseSummary, error) {
	var total float64
	if err := db.DB.Model(&models.Expense{}).Select("coalesce(sum(amount), 0)").Scan(&total).Error; err != nil {
		return nil, err
	}

	rows, err := db.DB.Model(&models.Expense{}).Select("category, sum(amount)").Group("category").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byCategory := make(map[string]float64)
	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err != nil {
			return nil, err
		}
		byCategory[category] = amount
	}

	return &ExpenseSummary{
		TotalExpenses: total,
		ByCategory:    byCategory,
	}, nil
}

func (s *ReportingService) GetPettyCashSummary() (*PettyCashSummary, error) {
	var credits float64
	var debits float64

	if err := db.DB.Model(&models.PettyCashTransaction{}).Where("type = ?", "credit").Select("coalesce(sum(amount), 0)").Scan(&credits).Error; err != nil {
		return nil, err
	}

	if err := db.DB.Model(&models.PettyCashTransaction{}).Where("type = ?", "debit").Select("coalesce(sum(amount), 0)").Scan(&debits).Error; err != nil {
		return nil, err
	}

	return &PettyCashSummary{
		TotalCredits: credits,
		TotalDebits:  debits,
		Balance:      credits - debits,
	}, nil
}
