CREATE TABLE daily_summaries (
    id              BIGSERIAL PRIMARY KEY,
    date            DATE NOT NULL UNIQUE,
    total_orders    INT NOT NULL DEFAULT 0,
    total_revenue   DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_customers INT NOT NULL DEFAULT 0,
    popular_item_id BIGINT REFERENCES menu_items(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_daily_summaries_date ON daily_summaries(date);