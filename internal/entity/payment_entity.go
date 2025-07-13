package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentOption string

const (
	PaymentOption5K   PaymentOption = "5k"   // 5k VND = 500 tokens
	PaymentOption10K  PaymentOption = "10k"  // 10k VND = 1000 tokens
	PaymentOption20K  PaymentOption = "20k"  // 20k VND = 2000 tokens
	PaymentOption50K  PaymentOption = "50k"  // 50k VND = 5000 tokens
	PaymentOption100K PaymentOption = "100k" // 100k VND = 10000 tokens
)

type PaymentTransaction struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        uint64             `json:"user_id" bson:"user_id"`
	TransactionID string             `json:"transaction_id" bson:"transaction_id"` // Auto-generated unique ID
	PaymentOption PaymentOption      `json:"payment_option" bson:"payment_option"`
	TokenAmount   int64              `json:"token_amount" bson:"token_amount"` // Amount in tokens
	VNDAmount     int64              `json:"vnd_amount" bson:"vnd_amount"`     // Amount in VND
	Status        PaymentStatus      `json:"status" bson:"status"`
	QRCode        string             `json:"qr_code" bson:"qr_code"`         // QR code text
	QRDataURL     string             `json:"qr_data_url" bson:"qr_data_url"` // Base64 image data
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
	CompletedAt   *time.Time         `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
}

// VietQRRequest represents the request to VietQR API
type VietQRRequest struct {
	AccountNo   string `json:"accountNo"`
	AccountName string `json:"accountName"`
	AcqID       string `json:"acqId"`
	Amount      int64  `json:"amount"`
	AddInfo     string `json:"addInfo"`
	Format      string `json:"format"`
	Template    string `json:"template"`
}

// VietQRResponse represents the response from VietQR API
type VietQRResponse struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
	Data struct {
		AcpID       string `json:"acpId"`
		AccountName string `json:"accountName"`
		QRCode      string `json:"qrCode"`
		QRDataURL   string `json:"qrDataURL"`
	} `json:"data"`
}

// PaymentOptionInfo contains payment option details
type PaymentOptionInfo struct {
	Option      PaymentOption `json:"option"`
	TokenAmount int64         `json:"token_amount"`
	VNDAmount   int64         `json:"vnd_amount"`
	Description string        `json:"description"`
}

func GetPaymentOptions() []PaymentOptionInfo {
	return []PaymentOptionInfo{
		{
			Option:      PaymentOption5K,
			TokenAmount: 500,
			VNDAmount:   5000,
			Description: "500 tokens - 5,000 VND",
		},
		{
			Option:      PaymentOption10K,
			TokenAmount: 1500,
			VNDAmount:   10000,
			Description: "1,500 tokens - 10,000 VND",
		},
		{
			Option:      PaymentOption20K,
			TokenAmount: 3500,
			VNDAmount:   20000,
			Description: "3,500 tokens - 20,000 VND",
		},
		{
			Option:      PaymentOption50K,
			TokenAmount: 10000,
			VNDAmount:   50000,
			Description: "10,000 tokens - 50,000 VND",
		},
		{
			Option:      PaymentOption100K,
			TokenAmount: 25000,
			VNDAmount:   100000,
			Description: "25,000 tokens - 100,000 VND",
		},
	}
}

func GetPaymentOptionInfo(option PaymentOption) *PaymentOptionInfo {
	options := GetPaymentOptions()
	for _, opt := range options {
		if opt.Option == option {
			return &opt
		}
	}
	return nil
}
