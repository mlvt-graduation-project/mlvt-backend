package entity

import "time"

type ClaimType string

const (
	ClaimDaily   ClaimType = "daily"
	ClaimPremium ClaimType = "premium"
)

// exactly mirrors daily_token_claims table
type TokenClaim struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	ClaimedDate time.Time `json:"claimed_date"`
	ClaimType   ClaimType `json:"claim_type"`
	Tokens      int64     `json:"tokens"`
	CreatedAt   time.Time `json:"created_at"`
}

// premium_users table
type PremiumUser struct {
	UserID    uint64    `json:"user_id"`
	ExpiredAt time.Time `json:"expired_at"`
}
