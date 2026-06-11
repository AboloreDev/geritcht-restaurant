CREATE TABLE refunds (
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT NOT NULL UNIQUE REFERENCES orders(id),
    payment_id      BIGINT NOT NULL REFERENCES payments(id),
    reference       VARCHAR(255) NOT NULL UNIQUE,
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,
    amount          DECIMAL(10,2) NOT NULL,
    currency        VARCHAR(10) DEFAULT 'NGN',
    status          VARCHAR(20) DEFAULT 'pending',
    reason          TEXT,
    processed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
