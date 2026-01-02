CREATE TYPE job_status AS ENUM ('pending', 'processing', 'completed', 'failed');

CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    api_key_id INTEGER references api_keys(id),
    url VARCHAR(2048) NOT NULL,
    status job_status DEFAULT 'pending',
    result_location VARCHAR(2048) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- ORM by default expects an index on deleted_at
CREATE INDEX idx_job_deleted_at ON jobs(deleted_at); 
