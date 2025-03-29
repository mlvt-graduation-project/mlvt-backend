package entity

import "time"

type WalletTransactionType string

const (
	TransactionTypeDeposit  WalletTransactionType = "deposit"
	TransactionTypeUseToken WalletTransactionType = "use_token"
)

type WalletTransaction struct {
	ID          uint64                `json:"id"`
	UserID      uint64                `json:"user_id"`
	Type        WalletTransactionType `json:"type"`         // deposit or use token
	Amount      uint64                `json:"amount"`       // token-based amount
	VoucherCode *string               `json:"voucher_code"` // optional voucher reference
	CreatedAt   time.Time             `json:"created_at"`
}
