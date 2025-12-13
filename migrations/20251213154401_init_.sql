-- +goose Up
-- +goose StatementBegin
-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
                                          id BIGSERIAL PRIMARY KEY,
                                          name VARCHAR(255) NOT NULL UNIQUE
);

-- Create nominants table
CREATE TABLE IF NOT EXISTS nominants (
                                         id BIGSERIAL PRIMARY KEY,
                                         name VARCHAR(255) NOT NULL,
                                         category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
                                     tg_id BIGINT PRIMARY KEY
);

-- Create votes table with unique constraint
CREATE TABLE IF NOT EXISTS votes (
                                     tg_user_id BIGINT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
                                     nominant_id BIGINT NOT NULL REFERENCES nominants(id) ON DELETE CASCADE,
                                     category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     UNIQUE(tg_user_id, category_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_nominants_category_id ON nominants(category_id);
CREATE INDEX IF NOT EXISTS idx_votes_tg_user_id ON votes(tg_user_id);
CREATE INDEX IF NOT EXISTS idx_votes_category_id ON votes(category_id);
CREATE INDEX IF NOT EXISTS idx_votes_nominant_id ON votes(nominant_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
