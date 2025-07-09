-- 1. Add the new column (default to 'daily' so existing rows keep working)
ALTER TABLE token_claims
  ADD COLUMN claim_type TEXT NOT NULL DEFAULT 'daily';

-- 2. Drop the old unique index on (user_id, claimed_date) if you created one
DROP INDEX IF EXISTS idx_token_claims_user_date;

-- 3. Create the new composite unique index
CREATE UNIQUE INDEX idx_token_claims_user_date_type
  ON token_claims(user_id, claimed_date, claim_type);
