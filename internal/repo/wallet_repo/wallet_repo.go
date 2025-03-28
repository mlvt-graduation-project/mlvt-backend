package wallet_repo

import (
	"context"
	"database/sql"
	"errors"
)

type WalletRepository interface {
	Deposit(ctx context.Context, userID uint64, amount int64) error
	Withdraw(ctx context.Context, userID uint64, amount int64) error
	GetBalance(ctx context.Context, userID uint64) (int64, error)
}

type walletRepo struct {
	db *sql.DB
}

func NewWalletRepo(db *sql.DB) WalletRepository {
	return &walletRepo{db: db}
}

func (r *walletRepo) Deposit(ctx context.Context, userID uint64, amount int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance + ? WHERE id = ?", amount, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *walletRepo) Withdraw(ctx context.Context, userID uint64, amount int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM users WHERE id = ?", userID).Scan(&currentBalance)
	if err != nil {
		return err
	}

	if currentBalance < amount {
		return errors.New("insufficient balance")
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance - ? WHERE id = ?", amount, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *walletRepo) GetBalance(ctx context.Context, userID uint64) (int64, error) {
	var balance int64
	err := r.db.QueryRowContext(ctx, "SELECT balance FROM users WHERE id = ?", userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
