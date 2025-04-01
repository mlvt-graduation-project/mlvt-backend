package wallet_service

import (
	"context"
	"errors"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/repo/wallet_repo"
	"mlvt/internal/service/traffic_service"
	"time"
)

type WalletService interface {
	Deposit(ctx context.Context, userID uint64, amount int64) error
	UseToken(ctx context.Context, userID uint64, amount int64) error
	GetBalance(ctx context.Context, userID uint64) (int64, error)
}

type walletService struct {
	repo           wallet_repo.WalletRepository
	trafficService traffic_service.TrafficService
}

func NewWalletService(
	repo wallet_repo.WalletRepository,
	trafficService traffic_service.TrafficService,
) WalletService {
	return &walletService{
		repo:           repo,
		trafficService: trafficService,
	}
}

func (s *walletService) Deposit(ctx context.Context, userID uint64, amount int64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	resultErr := s.repo.Deposit(ctx, userID, amount)

	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.DepositTokenAction,
		Description: fmt.Sprintf("account ID: %d, action deposit", userID),
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: deposit to user account, err: %s", err)
	}

	return resultErr
}

func (s *walletService) UseToken(ctx context.Context, userID uint64, amount int64) error {
	if amount <= 0 {
		return errors.New("to use tokens, the user's wallet amount must be positive")
	}
	resultErr := s.repo.UseToken(ctx, userID, amount)

	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.UseTokenAction,
		Description: fmt.Sprintf("account ID: %d, action use token", userID),
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: user use token, err: %s", err)
	}

	return resultErr

}

func (s *walletService) GetBalance(ctx context.Context, userID uint64) (int64, error) {
	amount, resultErr := s.repo.GetBalance(ctx, userID)

	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.DepositTokenAction,
		Description: fmt.Sprintf("get balance of account ID: %d", userID),
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: get balance of user account, err: %s", err)
	}

	return amount, resultErr
}
