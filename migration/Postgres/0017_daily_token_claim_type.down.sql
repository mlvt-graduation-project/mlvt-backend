-- 1. Drop the composite unique index
DROP INDEX IF EXISTS idx_token_claims_user_date_type;

-- 2. Drop the new column
ALTER TABLE token_claims DROP COLUMN IF EXISTS claim_type;

-- 3. Re-create the old unique index (on user_id, claimed_date)
CREATE UNIQUE INDEX idx_token_claims_user_date
  ON token_claims(user_id, claimed_date);