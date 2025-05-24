-- Add user_id column to books table
ALTER TABLE books
ADD COLUMN user_id BIGINT NOT NULL,
ADD FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE; 