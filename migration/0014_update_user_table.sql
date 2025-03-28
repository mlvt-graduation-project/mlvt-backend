PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;

-- 1) Rename the old table
ALTER TABLE users RENAME TO users_old;

-- 2) Create the new table with the updated schema
CREATE TABLE users (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name     TEXT,
    last_name      TEXT,
    username       TEXT,
    email          TEXT,
    password       TEXT,
    status         TEXT,            -- 'active' or 'inactive'
    role           TEXT,            -- 'user' or 'admin'
    avatar         TEXT,
    avatar_folder  TEXT,
    wallet_balance INTEGER,
    created_at     DATETIME,
    updated_at     DATETIME
);

-- 3) Copy/transform data from old table to the new one
--    - Old int status => new TEXT ('active'/'inactive')
--    - Old role => new role ('admin'/'user'), 
--        assuming the old role was stored as exactly "Admin" or "User"
--    - Old "premium" is removed
--    - Initialize wallet_balance to 0 (or another value if you'd like)
INSERT INTO users (
    id, first_name, last_name, username, email, password, status, role,
    avatar, avatar_folder, wallet_balance, created_at, updated_at
)
SELECT
    id,
    first_name,
    last_name,
    username,
    email,
    password,

    -- Convert your old integer-based status to new "active"/"inactive"
    CASE 
        WHEN status = 1 THEN 'active'
        ELSE 'inactive'
    END AS status,

    -- Map old role string to your new user/admin. Adjust as needed.
    CASE 
        WHEN role = 'Admin' THEN 'admin'
        ELSE 'user'
    END AS role,

    avatar,
    avatar_folder,

    -- Initialize wallet_balance. You can also dynamically set it based on old 'premium' if needed.
    0 AS wallet_balance,

    created_at,
    updated_at
FROM users_old;

-- 4) Drop the old table
DROP TABLE users_old;

COMMIT;
PRAGMA foreign_keys = ON;
