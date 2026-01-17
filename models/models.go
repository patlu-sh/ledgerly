package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "credit"
	TransactionTypeDebit  TransactionType = "debit"
)

type PettyCashTransaction struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Description string          `json:"description"`
	UserID      string          `json:"user_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

type Expense struct {
	ID                      uint           `gorm:"primaryKey" json:"id"`
	Title                   string         `json:"title"`
	Amount                  float64        `json:"amount"`
	Category                string         `json:"category"`
	UserID                  string         `json:"user_id"`
	PettyCashTransactionID  *uint          `json:"petty_cash_transaction_id"`
	PettyCashTransaction    *PettyCashTransaction `gorm:"foreignKey:PettyCashTransactionID" json:"petty_cash_transaction,omitempty"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`
}
