package payment_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/repo/payment_repo"
	"mlvt/internal/service/traffic_service"
	"mlvt/internal/service/wallet_service"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, userID uint64, option entity.PaymentOption) (*entity.PaymentTransaction, error)
	GetPaymentByID(ctx context.Context, paymentID primitive.ObjectID) (*entity.PaymentTransaction, error)
	GetPaymentByTransactionID(ctx context.Context, transactionID string) (*entity.PaymentTransaction, error)
	GetUserPayments(ctx context.Context, userID uint64) ([]entity.PaymentTransaction, error)
	ConfirmPayment(ctx context.Context, transactionID string) error
	CancelPayment(ctx context.Context, paymentID primitive.ObjectID) error
	GetPaymentOptions() []entity.PaymentOptionInfo
	GetPendingPayments(ctx context.Context) ([]entity.PaymentTransaction, error)
}

type paymentService struct {
	paymentRepo    payment_repo.PaymentRepository
	walletService  wallet_service.WalletService
	trafficService traffic_service.TrafficService
}

func NewPaymentService(
	paymentRepo payment_repo.PaymentRepository,
	walletService wallet_service.WalletService,
	trafficService traffic_service.TrafficService,
) PaymentService {
	return &paymentService{
		paymentRepo:    paymentRepo,
		walletService:  walletService,
		trafficService: trafficService,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, userID uint64, option entity.PaymentOption) (*entity.PaymentTransaction, error) {
	// Get payment option info
	optionInfo := entity.GetPaymentOptionInfo(option)
	if optionInfo == nil {
		return nil, fmt.Errorf("invalid payment option: %s", option)
	}

	// Generate unique transaction ID
	transactionID := generateTransactionID()

	// Create payment transaction
	payment := &entity.PaymentTransaction{
		UserID:        userID,
		TransactionID: transactionID,
		PaymentOption: option,
		TokenAmount:   optionInfo.TokenAmount,
		VNDAmount:     optionInfo.VNDAmount,
		Status:        entity.PaymentStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Call VietQR API to generate QR code
	qrCode, qrDataURL, err := s.callVietQRAPI(transactionID, optionInfo.VNDAmount)
	if err != nil {
		log.Errorf("Failed to generate QR code: %v", err)
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	payment.QRCode = qrCode
	payment.QRDataURL = qrDataURL

	// Save to database
	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		log.Errorf("Failed to create payment: %v", err)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Log traffic
	s.logTraffic(ctx, userID, "create_payment", fmt.Sprintf("Created payment %s for %d tokens", transactionID, optionInfo.TokenAmount))

	return payment, nil
}

func (s *paymentService) GetPaymentByID(ctx context.Context, paymentID primitive.ObjectID) (*entity.PaymentTransaction, error) {
	return s.paymentRepo.GetByID(ctx, paymentID)
}

func (s *paymentService) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*entity.PaymentTransaction, error) {
	return s.paymentRepo.GetByTransactionID(ctx, transactionID)
}

func (s *paymentService) GetUserPayments(ctx context.Context, userID uint64) ([]entity.PaymentTransaction, error) {
	return s.paymentRepo.GetByUserID(ctx, userID)
}

func (s *paymentService) ConfirmPayment(ctx context.Context, transactionID string) error {
	// Get payment by transaction ID
	payment, err := s.paymentRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return fmt.Errorf("payment not found")
	}

	// Check if payment is still pending
	if payment.Status != entity.PaymentStatusPending {
		return fmt.Errorf("payment is not pending")
	}

	// Mark payment as completed
	err = s.paymentRepo.MarkAsCompleted(ctx, payment.ID)
	if err != nil {
		return fmt.Errorf("failed to mark payment as completed: %w", err)
	}

	// Add tokens to user wallet
	err = s.walletService.Deposit(ctx, payment.UserID, payment.TokenAmount)
	if err != nil {
		log.Errorf("Failed to deposit tokens to user wallet: %v", err)
		// Revert payment status
		s.paymentRepo.UpdateStatus(ctx, payment.ID, entity.PaymentStatusFailed)
		return fmt.Errorf("failed to deposit tokens: %w", err)
	}

	// Log traffic
	s.logTraffic(ctx, payment.UserID, "confirm_payment", fmt.Sprintf("Confirmed payment %s, added %d tokens", transactionID, payment.TokenAmount))

	return nil
}

func (s *paymentService) CancelPayment(ctx context.Context, paymentID primitive.ObjectID) error {
	return s.paymentRepo.UpdateStatus(ctx, paymentID, entity.PaymentStatusCancelled)
}

func (s *paymentService) GetPaymentOptions() []entity.PaymentOptionInfo {
	return entity.GetPaymentOptions()
}

func (s *paymentService) GetPendingPayments(ctx context.Context) ([]entity.PaymentTransaction, error) {
	return s.paymentRepo.GetPendingPayments(ctx)
}

func (s *paymentService) callVietQRAPI(transactionID string, amount int64) (string, string, error) {
	// Prepare request
	vietQRRequest := entity.VietQRRequest{
		AccountNo:   env.EnvConfig.VietinBankAccountNo,
		AccountName: env.EnvConfig.VietinBankAccountName,
		AcqID:       env.EnvConfig.VietinBankBinCode,
		Amount:      amount,
		AddInfo:     fmt.Sprintf("MLVT-%s", transactionID),
		Format:      "text",
		Template:    "compact",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(vietQRRequest)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.vietqr.io/v2/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-client-id", env.EnvConfig.VietQRClientID)
	req.Header.Set("x-api-key", env.EnvConfig.VietQRAPIKey)

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var vietQRResponse entity.VietQRResponse
	if err := json.NewDecoder(resp.Body).Decode(&vietQRResponse); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check response code
	if vietQRResponse.Code != "00" {
		return "", "", fmt.Errorf("VietQR API error: %s - %s", vietQRResponse.Code, vietQRResponse.Desc)
	}

	return vietQRResponse.Data.QRCode, vietQRResponse.Data.QRDataURL, nil
}

func (s *paymentService) logTraffic(ctx context.Context, userID uint64, action, description string) {
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.TrafficActionType(action),
		Description: fmt.Sprintf("User %d: %s", userID, description),
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("Failed to log traffic: %v", err)
	}
}

func generateTransactionID() string {
	return fmt.Sprintf("TXN_%d", time.Now().UnixNano())
}
