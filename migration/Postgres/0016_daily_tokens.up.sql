-- who claimed what & when
CREATE TABLE token_claims (
  id            SERIAL PRIMARY KEY,
  user_id       INTEGER NOT NULL,
  claimed_date  DATE    NOT NULL,
  tokens        INTEGER NOT NULL,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  claim_type    TEXT NOT NULL DEFAULT 'daily',
  UNIQUE(user_id, claimed_date, claim_type)
);

-- premium flag with expiry
CREATE TABLE premium_users (
    user_id    BIGINT PRIMARY KEY,
    expired_at TIMESTAMP NOT NULL                -- when the party’s over
);


CREATE UNIQUE INDEX idx_token_claims_user_date_type ON token_claims(user_id, claimed_date, claim_type);
DROP INDEX IF EXISTS idx_token_claims_user_date;