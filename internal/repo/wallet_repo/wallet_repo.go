package wallet_repo

import (
	"context"
	"errors"
	"mlvt/internal/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	Deposit(ctx context.Context, userID uint64, amount int64) error
	UseToken(ctx context.Context, userID uint64, amount int64) error
	GetBalance(ctx context.Context, userID uint64) (int64, error)
}

type walletRepo struct {
	db *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) WalletRepository {
	return &walletRepo{db: db}
}

func (r *walletRepo) Deposit(ctx context.Context, userID uint64, amount int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateQuery := r.db.Rebind("UPDATE users SET wallet_balance = wallet_balance + ? WHERE id = ?")
	_, err = tx.ExecContext(ctx, updateQuery, amount, userID)
	if err != nil {
		return err
	}

	insertQuery := r.db.Rebind(`
		INSERT INTO wallet_transactions (user_id, type, amount, created_at)
		VALUES (?, ?, ?, ?)`)
	_, err = tx.ExecContext(ctx, insertQuery,
		userID,
		entity.TransactionTypeDeposit,
		amount,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *walletRepo) UseToken(ctx context.Context, userID uint64, amount int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	selectQuery := r.db.Rebind("SELECT wallet_balance FROM users WHERE id = ?")
	var currentBalance int64
	err = tx.QueryRowContext(ctx, selectQuery, userID).Scan(&currentBalance)
	if err != nil {
		return err
	}

	if currentBalance < amount {
		return errors.New("insufficient balance")
	}

	updateQuery := r.db.Rebind("UPDATE users SET wallet_balance = wallet_balance - ? WHERE id = ?")
	_, err = tx.ExecContext(ctx, updateQuery, amount, userID)
	if err != nil {
		return err
	}

	insertQuery := r.db.Rebind(`
		INSERT INTO wallet_transactions (user_id, type, amount, created_at)
		VALUES (?, ?, ?, ?)`)
	_, err = tx.ExecContext(ctx, insertQuery,
		userID,
		entity.TransactionTypeUseToken,
		amount,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *walletRepo) GetBalance(ctx context.Context, userID uint64) (int64, error) {
	query := r.db.Rebind("SELECT wallet_balance FROM users WHERE id = ?")
	var balance int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
