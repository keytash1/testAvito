package repository

import (
	"pr_reviewer_service_go/internal/db"
	"pr_reviewer_service_go/internal/repository/interfaces"

	"gorm.io/gorm"
)

type TransactionRepository struct{}

func NewTransactionRepository() interfaces.TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}
