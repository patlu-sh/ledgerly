package services

import (
	"errors"
	"ledgerly/db"
	"ledgerly/models"
)

type PettyCashService struct{}

func (s *PettyCashService) CreateTransaction(tx *models.PettyCashTransaction) error {
	if tx.Type == models.TransactionTypeDebit {
		balance, err := s.GetBalance()
		if err != nil {
			return err
		}
		if balance < tx.Amount {
			return errors.New("insufficient funds")
		}
	}
	return db.DB.Create(tx).Error
}

func (s *PettyCashService) GetBalance() (float64, error) {
	var credits float64
	var debits float64

	if err := db.DB.Model(&models.PettyCashTransaction{}).Where("type = ?", models.TransactionTypeCredit).Select("coalesce(sum(amount), 0)").Scan(&credits).Error; err != nil {
		return 0, err
	}

	if err := db.DB.Model(&models.PettyCashTransaction{}).Where("type = ?", models.TransactionTypeDebit).Select("coalesce(sum(amount), 0)").Scan(&debits).Error; err != nil {
		return 0, err
	}

	return credits - debits, nil
}

func (s *PettyCashService) ListTransactions() ([]models.PettyCashTransaction, error) {
	var transactions []models.PettyCashTransaction
	err := db.DB.Find(&transactions).Error
	return transactions, err
}
