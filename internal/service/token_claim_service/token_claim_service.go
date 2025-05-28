package token_claim_service

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/repo/token_claim_repo"
)

type TokenService interface {
	ClaimDaily(ctx context.Context, userID uint64) error   // +5
	ClaimPremium(ctx context.Context, userID uint64) error // +20
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
	return s.repo.Claim(ctx, userID, 5)
}

func (s *svc) ClaimPremium(ctx context.Context, userID uint64) error {
	ok, err := s.repo.IsPremium(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return token_claim_repo.ErrNotPremium
	}
	return s.repo.Claim(ctx, userID, 20)
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
