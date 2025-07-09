-- Xóa foreign key từ videos.user_id → users.id
ALTER TABLE videos DROP CONSTRAINT IF EXISTS videos_user_id_fkey;

-- Xóa foreign key từ transcriptions.user_id → users.id
ALTER TABLE transcriptions DROP CONSTRAINT IF EXISTS transcriptions_user_id_fkey;

-- Xóa foreign key từ audios.user_id → users.id
ALTER TABLE audios DROP CONSTRAINT IF EXISTS audios_user_id_fkey;

DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    status INTEGER NOT NULL,
    premium BOOLEAN NOT NULL DEFAULT FALSE,
    role TEXT NOT NULL DEFAULT 'User',
    avatar TEXT NOT NULL DEFAULT '',
    avatar_folder TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE videos ADD CONSTRAINT videos_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE transcriptions ADD CONSTRAINT transtiptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE audios ADD CONSTRAINT audios_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;