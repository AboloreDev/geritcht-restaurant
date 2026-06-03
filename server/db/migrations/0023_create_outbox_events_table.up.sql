CREATE TABLE outbox_events (
    id           BIGSERIAL PRIMARY KEY,
    event_type   VARCHAR(100) NOT NULL,
    payload      JSONB NOT NULL,
    status       VARCHAR(20) DEFAULT 'pending',
    retry_count  INT DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_status_created ON outbox_events (status, created_at);