CREATE TYPE table_status AS ENUM ('available', 'occupied', 'reserved');

CREATE TABLE tables (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    capacity   INT NOT NULL,
    location   VARCHAR(100),
    status     table_status NOT NULL DEFAULT 'available',
    qr_code    VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- add foreign keys now that tables exists
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_table_id
    FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE SET NULL;

ALTER TABLE reservations
    ADD CONSTRAINT fk_reservations_table_id
    FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE RESTRICT;

CREATE INDEX idx_tables_status ON tables(status);
CREATE INDEX idx_tables_deleted_at ON tables(deleted_at);
CREATE INDEX idx_tables_id ON tables(id);
CREATE INDEX idx_tables_name ON tables(name);
CREATE INDEX idx_tables_capacity ON tables(capacity);
CREATE INDEX idx_tables_id_capacity ON tables(id, capacity);