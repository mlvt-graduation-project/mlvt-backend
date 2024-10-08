package repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
)

// TransactionLogRepo is responsible for logging transaction events to the database
type TransactionLogRepo interface {
	LogTransaction(log *entity.TransactionLog) error
}

type transactionLogRepo struct {
	db *sql.DB
}

// NewTransactionLogRepo creates a new repository for logging transactions
func NewTransactionLogRepo(db *sql.DB) TransactionLogRepo {
	return &transactionLogRepo{db: db}
}

// LogTransaction inserts a log entry into the transaction_logs table
func (r *transactionLogRepo) LogTransaction(log *entity.TransactionLog) error {
	query := `INSERT INTO transaction_logs (order_id, payment_method, action, status, details)
              VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, log.OrderID, log.PaymentMethod, log.Action, log.Status, log.Details)
	if err != nil {
		fmt.Errorf("Error logging transaction: %v", err)
		return err
	}
	return nil
}
