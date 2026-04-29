CREATE TABLE order_items (
    id                   BIGSERIAL PRIMARY KEY,
    order_id             BIGINT NOT NULL,
    menu_item_id         BIGINT NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
    quantity             INT NOT NULL,
    price                DECIMAL(10,2) NOT NULL,
    special_instructions TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_menu_item_id ON order_items(menu_item_id);
CREATE INDEX idx_order_items_deleted_at ON order_items(deleted_at);