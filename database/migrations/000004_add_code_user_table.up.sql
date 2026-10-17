ALTER TABLE users ADD COLUMN code VARCHAR(255) NULL;

CREATE UNIQUE INDEX users_code_unique_idx ON users (code);