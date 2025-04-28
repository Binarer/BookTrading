-- Add states table
CREATE TABLE IF NOT EXISTS states (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert default states
INSERT INTO states (name) VALUES 
    ('available'),
    ('trading'),
    ('traded');

-- Add foreign key to books table
ALTER TABLE books
ADD COLUMN state_id BIGINT,
ADD FOREIGN KEY (state_id) REFERENCES states(id); 