-- who claimed what & when
CREATE TABLE token_claims (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id       INTEGER NOT NULL,
  claimed_date  DATE    NOT NULL,
  tokens        INTEGER NOT NULL,
  created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, claimed_date)
);

-- premium flag with expiry
CREATE TABLE premium_users (
    user_id    BIGINT PRIMARY KEY,
    expired_at TIMESTAMP NOT NULL                -- when the party’s over
);
