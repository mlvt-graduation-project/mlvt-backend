INSERT INTO users (first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at) VALUES
('John', 'Doe', 'johndoe', 'john@example.com', 'hashedpassword1', 1, FALSE, 'User', 'avatar1.png', 'avatars/', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Jane', 'Smith', 'janesmith', 'jane@example.com', 'hashedpassword2', 1, TRUE, 'Admin', 'avatar2.png', 'avatars/', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT(email) DO NOTHING;
