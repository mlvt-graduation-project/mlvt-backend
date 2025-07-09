package token_claim_repo

import (
	"context"
	"database/sql"
	"errors"
	"mlvt/internal/entity"
	"time"

	"github.com/jmoiron/sqlx"
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
	GetLastClaimDate(ctx context.Context, userID uint64, ctype entity.ClaimType) (time.Time, error)
	ClaimAtDate(ctx context.Context, userID uint64, amount int64, ctype entity.ClaimType, date time.Time) error
}

type tokenRepo struct{ db *sqlx.DB }

func New(db *sqlx.DB) TokenRepository { return &tokenRepo{db} }

// atomic: insert claim → credit wallet → write wallet_tx
func (r *tokenRepo) Claim(ctx context.Context, userID uint64, amount int64, ctype entity.ClaimType) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. record claim by type
	query1 := `
		INSERT INTO token_claims
		(user_id, claimed_date, claim_type, tokens)
		VALUES
		(?, CURRENT_DATE, ?, ?)
		ON CONFLICT(user_id, claimed_date, claim_type) DO NOTHING`
	query1 = sqlx.Rebind(sqlx.DOLLAR, query1)

	res, err := tx.ExecContext(ctx, query1, userID, ctype, amount)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrAlreadyClaimed
	}

	// 2. credit wallet
	query2 := `
		UPDATE users 
		SET wallet_balance = wallet_balance + ?
		WHERE id = ?`
	query2 = sqlx.Rebind(sqlx.DOLLAR, query2)

	if _, err := tx.ExecContext(ctx, query2, amount, userID); err != nil {
		return err
	}

	// 3. log transaction
	query3 := `
		INSERT INTO wallet_transactions
		(user_id, type, amount, created_at)
		VALUES (?, ?, ?, NOW())`
	query3 = sqlx.Rebind(sqlx.DOLLAR, query3)

	if _, err := tx.ExecContext(ctx, query3, userID, entity.TransactionTypeDeposit, amount); err != nil {
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
	// 1) fetch the old expiry (if any)
	var old sql.NullTime
	query1 := `SELECT expired_at FROM premium_users WHERE user_id = ?`
	query1 = sqlx.Rebind(sqlx.DOLLAR, query1)

	err := r.db.QueryRowContext(ctx, query1, userID).Scan(&old)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	now := time.Now()
	isExtension := old.Valid && old.Time.After(now)

	// 2) compute new expiry
	base := now
	if isExtension {
		base = old.Time
	}
	newExp := base.AddDate(0, 0, 30)

	// 3) UPSERT
	query2 := `
		INSERT INTO premium_users (user_id, expired_at)
		VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE
		  SET expired_at = excluded.expired_at`
	query2 = sqlx.Rebind(sqlx.DOLLAR, query2)

	if _, err := r.db.ExecContext(ctx, query2, userID, newExp); err != nil {
		return err
	}

	// 4) Grant bonus if this is NOT an extension
	if !isExtension {
		if err := r.Claim(ctx, userID, 20, entity.ClaimPremium); err != nil {
			return err
		}
	}

	return nil
}

func (r *tokenRepo) IsPremium(ctx context.Context, userID uint64) (bool, error) {
	if err := r.purgeExpired(ctx); err != nil {
		return false, err
	}

	var exp time.Time
	query := `SELECT expired_at FROM premium_users WHERE user_id = ?`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exp)
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

func (r *tokenRepo) GetLastClaimDate(ctx context.Context, userID uint64, ctype entity.ClaimType) (time.Time, error) {
	query := `
		SELECT MAX(claimed_date) FROM token_claims 
		WHERE user_id = ? AND claim_type = ?`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var dt sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID, ctype).Scan(&dt)
	if err != nil {
		return time.Time{}, err
	}
	if !dt.Valid {
		return time.Time{}, nil // zero → never claimed
	}

	t, _ := time.Parse("2006-01-02", dt.String)
	return t, nil
}

// 2) Backfill a specific date
func (r *tokenRepo) ClaimAtDate(ctx context.Context, userID uint64, amount int64, ctype entity.ClaimType, date time.Time) error {
	query := `
		INSERT INTO token_claims (user_id, claimed_date, claim_type, tokens)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, claimed_date, claim_type) DO NOTHING`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	_, err := r.db.ExecContext(ctx,
		query,
		userID,
		date.Format("2006-01-02"),
		ctype,
		amount,
	)
	return err
}
