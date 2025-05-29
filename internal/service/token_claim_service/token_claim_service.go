package token_claim_service

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/repo/token_claim_repo"
	"time"
)

type TokenService interface {
	ClaimDaily(ctx context.Context, userID uint64) error                          // +5
	ClaimPremium(ctx context.Context, userID uint64) (daysClaimed int, err error) // +20
	ListClaims(ctx context.Context) ([]entity.TokenClaim, error)
	AddPremium(ctx context.Context, userID uint64) error
	ListPremiumUsers(ctx context.Context) ([]entity.PremiumUser, error)
	CheckPremium(ctx context.Context, userID uint64) (bool, error)
}

type svc struct {
	repo token_claim_repo.TokenRepository
}

func New(r token_claim_repo.TokenRepository) TokenService { return &svc{r} }

func (s *svc) ClaimDaily(ctx context.Context, userID uint64) error {
	return s.repo.Claim(ctx, userID, 5, entity.ClaimDaily)
}

func (s *svc) ClaimPremium(ctx context.Context, userID uint64) (daysClaimed int, err error) {
	// 1) must still be active premium
	ok, err := s.repo.IsPremium(ctx, userID)
	if err != nil || !ok {
		return 0, token_claim_repo.ErrNotPremium
	}

	// 2) when did they last snag premium tokens?
	lastDate, err := s.repo.GetLastClaimDate(ctx, userID, entity.ClaimPremium)
	if err != nil {
		return 0, err
	}

	// 3) how many full days since then (excluding lastDate)?
	today := time.Now().UTC().Truncate(24 * time.Hour)
	var startDate time.Time
	if lastDate.IsZero() {
		// never claimed → start with yesterday so they get today's 20
		startDate = today.Add(-24 * time.Hour)
	} else {
		startDate = lastDate
	}
	days := int(today.Sub(startDate).Hours() / 24)
	if days <= 0 {
		return 0, token_claim_repo.ErrAlreadyClaimed
	}

	// 4) backfill each missing day
	for i := 1; i <= days; i++ {
		d := startDate.Add(time.Duration(i) * 24 * time.Hour)
		if err := s.repo.ClaimAtDate(ctx, userID, 20, entity.ClaimPremium, d); err != nil {
			return 0, err
		}
	}

	return days, nil
}

func (s *svc) AddPremium(ctx context.Context, userID uint64) error {
	return s.repo.AddPremium(ctx, userID)
}

func (s *svc) ListClaims(ctx context.Context) ([]entity.TokenClaim, error) {
	return s.repo.ListClaims(ctx)
}

func (s *svc) ListPremiumUsers(ctx context.Context) ([]entity.PremiumUser, error) {
	return s.repo.ListPremium(ctx)
}

func (s *svc) CheckPremium(ctx context.Context, userID uint64) (bool, error) {
	return s.repo.IsPremium(ctx, userID)
}
