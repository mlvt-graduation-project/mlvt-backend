package token_claim_repo

import (
	"context"
	"database/sql"
	"errors"
	"mlvt/internal/entity"
	"time"
)

var (
	ErrAlreadyClaimed = errors.New("token already claimed today")
	ErrNotPremium     = errors.New("user is not premium or expired")
)

type TokenRepository interface {
	Claim(ctx context.Context, userID uint64, amount int64, ctype entity.ClaimType) error
	IsPremium(ctx context.Context, userID uint64) (bool, error)
	AddPremium(ctx context.Context, userID uint64) error
	ListClaims(ctx context.Context) ([]entity.TokenClaim, error)
	ListPremium(ctx context.Context) ([]entity.PremiumUser, error)
}

type tokenRepo struct{ db *sql.DB }

func New(db *sql.DB) TokenRepository { return &tokenRepo{db} }

// atomic: insert claim → credit wallet → write wallet_tx
func (r *tokenRepo) Claim(ctx context.Context, userID uint64, amount int64, ctype entity.ClaimType) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. record claim by type
	res, err := tx.ExecContext(ctx,
		`INSERT INTO token_claims
       (user_id, claimed_date, claim_type, tokens)
     VALUES
       (?, date('now'), ?, ?)
     ON CONFLICT(user_id, claimed_date, claim_type) DO NOTHING`,
		userID, ctype, amount,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrAlreadyClaimed
	}

	// 2. credit wallet
	if _, err := tx.ExecContext(ctx,
		`UPDATE users 
        SET wallet_balance = wallet_balance + ? 
      WHERE id = ?`,
		amount, userID,
	); err != nil {
		return err
	}

	// 3. log transaction
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO wallet_transactions
       (user_id, type, amount, created_at)
     VALUES
       (?, ?, ?, datetime('now'))`,
		userID, entity.TransactionTypeDeposit, amount,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *tokenRepo) purgeExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM premium_users WHERE expired_at <= CURRENT_TIMESTAMP`,
	)
	return err
}

func (r *tokenRepo) AddPremium(ctx context.Context, userID uint64) error {
	// get current expiry (if any)
	var old sql.NullTime
	err := r.db.QueryRowContext(ctx,
		`SELECT expired_at FROM premium_users WHERE user_id = ?`, userID,
	).Scan(&old)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	base := time.Now()
	if old.Valid && old.Time.After(base) {
		base = old.Time
	}
	newExp := base.AddDate(0, 0, 30)

	// upsert
	_, err = r.db.ExecContext(ctx, `
        INSERT INTO premium_users (user_id, expired_at)
        VALUES (?, ?)
        ON CONFLICT(user_id) DO UPDATE
          SET expired_at = excluded.expired_at
    `, userID, newExp)
	return err
}

func (r *tokenRepo) IsPremium(ctx context.Context, userID uint64) (bool, error) {
	if err := r.purgeExpired(ctx); err != nil {
		return false, err
	}

	var exp time.Time
	err := r.db.QueryRowContext(ctx,
		`SELECT expired_at FROM premium_users WHERE user_id = ?`, userID).Scan(&exp)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return exp.After(time.Now()), nil
}

func (r *tokenRepo) ListClaims(ctx context.Context) ([]entity.TokenClaim, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, claimed_date, tokens, created_at
           FROM token_claims
         ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []entity.TokenClaim
	for rows.Next() {
		var c entity.TokenClaim
		if err := rows.Scan(&c.ID, &c.UserID, &c.ClaimedDate, &c.Tokens, &c.CreatedAt); err != nil {
			return nil, err
		}
		claims = append(claims, c)
	}
	return claims, rows.Err()
}

func (r *tokenRepo) ListPremium(ctx context.Context) ([]entity.PremiumUser, error) {
	if err := r.purgeExpired(ctx); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT user_id, expired_at
		   FROM premium_users
		  WHERE expired_at > CURRENT_TIMESTAMP
		  ORDER BY expired_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []entity.PremiumUser
	for rows.Next() {
		var p entity.PremiumUser
		if err := rows.Scan(&p.UserID, &p.ExpiredAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, rows.Err()
}
