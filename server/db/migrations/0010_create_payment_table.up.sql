CREATE TABLE payments (
    id                 BIGSERIAL PRIMARY KEY,
    order_id           BIGINT NOT NULL UNIQUE REFERENCES orders(id) ON DELETE RESTRICT,
    user_id            BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    reference          VARCHAR(255) NOT NULL UNIQUE,
    idempotency_key    VARCHAR(255) NOT NULL UNIQUE,
    amount             DECIMAL(10,2) NOT NULL,
    currency           VARCHAR(10) NOT NULL DEFAULT 'NGN',
    status             payment_status NOT NULL DEFAULT 'unpaid',
    provider           VARCHAR(50) NOT NULL DEFAULT 'paystack',
    provider_reference VARCHAR(255),
    failure_reason     TEXT,
    paid_at            TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_reference ON payments(reference);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_deleted_at ON payments(deleted_at);