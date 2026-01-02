CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL,
    hashed_secret TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_public_id ON api_keys(public_id); 
-- ORM by default expects an index on deleted_at
CREATE INDEX idx_api_key_deleted_at ON api_keys(deleted_at); 
