CREATE TYPE stock_movement_type AS ENUM ('in', 'out', 'waste');

CREATE TABLE stock_movements (
    id            BIGSERIAL PRIMARY KEY,
    ingredient_id BIGINT NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    type          stock_movement_type NOT NULL,
    quantity      DECIMAL(10,2) NOT NULL,
    reason        TEXT,
    created_by    BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stock_movements_ingredient_id ON stock_movements(ingredient_id);
CREATE INDEX idx_stock_movements_created_by ON stock_movements(created_by);
CREATE INDEX idx_stock_movements_created_at ON stock_movements(created_at);