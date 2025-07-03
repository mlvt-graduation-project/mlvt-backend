BEGIN;

-- Xóa foreign key từ videos.user_id → users.id
ALTER TABLE videos DROP CONSTRAINT IF EXISTS videos_user_id_fkey;

-- Xóa foreign key từ transcriptions.user_id → users.id
ALTER TABLE transcriptions DROP CONSTRAINT IF EXISTS transcriptions_user_id_fkey;

-- Xóa foreign key từ audios.user_id → users.id
ALTER TABLE audios DROP CONSTRAINT IF EXISTS audios_user_id_fkey;



-- Bắt đầu quá trình migration
-- 1) Đổi tên bảng cũ
ALTER TABLE users RENAME TO users_old;

-- 2) Tạo bảng mới
CREATE TABLE users (
    id             SERIAL PRIMARY KEY,
    first_name     TEXT,
    last_name      TEXT,
    username       TEXT,
    email          TEXT,
    password       TEXT,
    status         TEXT,
    role           TEXT,
    avatar         TEXT,
    avatar_folder  TEXT,
    wallet_balance INTEGER,
    created_at     TIMESTAMP,
    updated_at     TIMESTAMP
);

-- 3) Insert dữ liệu có transform
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
    CASE WHEN status = 1 THEN 'active' ELSE 'inactive' END,
    CASE WHEN role = 'Admin' THEN 'admin' ELSE 'user' END,
    avatar,
    avatar_folder,
    0,
    created_at,
    updated_at
FROM users_old;

-- 4) Xoá bảng cũ
DROP TABLE users_old;

-- 5) Thêm lại foreign key
ALTER TABLE videos ADD CONSTRAINT videos_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE transcriptions ADD CONSTRAINT transtiptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE audios ADD CONSTRAINT audios_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

COMMIT;
