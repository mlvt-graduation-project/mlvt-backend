package wallet_service

import (
	"context"
	"errors"
	"mlvt/internal/repo/wallet_repo"
)

type WalletService interface {
	Deposit(ctx context.Context, userID uint64, amount int64) error
	Withdraw(ctx context.Context, userID uint64, amount int64) error
	GetBalance(ctx context.Context, userID uint64) (int64, error)
}

type walletService struct {
	repo wallet_repo.WalletRepository
}

func NewWalletService(repo wallet_repo.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) Deposit(ctx context.Context, userID uint64, amount int64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	return s.repo.Deposit(ctx, userID, amount)
}

func (s *walletService) Withdraw(ctx context.Context, userID uint64, amount int64) error {
	if amount <= 0 {
		return errors.New("withdraw amount must be positive")
	}
	return s.repo.Withdraw(ctx, userID, amount)
}

func (s *walletService) GetBalance(ctx context.Context, userID uint64) (int64, error) {
	return s.repo.GetBalance(ctx, userID)
}
