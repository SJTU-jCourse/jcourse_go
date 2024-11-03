package repository

import (
	"context"
	"gorm.io/gorm"
)

type DBOperation func(*gorm.DB) error
type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) TransactionHandler {
	return TransactionHandler{db: db}
}

func InTransAction(ctx context.Context, q TransactionHandler, operation DBOperation) error {
	db := q.db.WithContext(ctx)
	tx := db.Begin()
	if err := operation(tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
