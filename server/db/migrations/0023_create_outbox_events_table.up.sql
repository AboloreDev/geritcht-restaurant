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
CREATE INDEX idx_outbox_status_retry_count ON outbox_events (status, retry_count);
CREATE INDEX idx_outbox_id_status_event_type ON outbox_events (id, status, event_type);
CREATE INDEX idx_outbox_event_type ON outbox_events (event_type);
CREATE INDEX idx_outbox_processed_at ON outbox_events (processed_at);
CREATE INDEX idx_outbox_retry_count ON outbox_events (retry_count);
CREATE INDEX idx_outbox_created_at ON outbox_events (created_at);
