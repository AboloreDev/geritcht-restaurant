CREATE TABLE recommendations (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT REFERENCES users(id) ON DELETE SET NULL,
    allergies     TEXT,
    dietary       VARCHAR(100),
    mood          VARCHAR(100),
    budget        VARCHAR(50),
    suggested_ids TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_recommendations_user_id ON recommendations(user_id);
CREATE INDEX idx_recommendations_deleted_at ON recommendations(deleted_at);