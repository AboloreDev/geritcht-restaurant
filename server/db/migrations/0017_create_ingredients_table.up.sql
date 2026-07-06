CREATE TABLE ingredients (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL UNIQUE,
    unit          VARCHAR(50) NOT NULL,
    current_stock DECIMAL(10,2) NOT NULL DEFAULT 0,
    min_threshold DECIMAL(10,2) NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_ingredients_deleted_at ON ingredients(deleted_at);
CREATE INDEX idx_ingredients_name ON ingredients(name);
CREATE INDEX idx_ingredients_current_stock_min_threshold ON ingredients(current_stock, min_threshold);
CREATE INDEX idx_ingredients_id ON ingredients(id);
CREATE INDEX idx_ingredients_created_at ON ingredients(created_at);
CREATE INDEX idx_ingredients_current_stock ON ingredients(current_stock);
CREATE INDEX idx_ingredients_min_threshold ON ingredients(min_threshold);