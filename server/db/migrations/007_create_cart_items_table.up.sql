CREATE TABLE cart_items (
    id                   BIGSERIAL PRIMARY KEY,
    cart_id              BIGINT NOT NULL REFERENCES cart(id) ON DELETE CASCADE,
    menu_id         BIGINT NOT NULL REFERENCES menu(id) ON DELETE RESTRICT,
    quantity             INT NOT NULL DEFAULT 1,
    special_instructions TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ
);

CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_menu_id ON cart_items(menu_id);
CREATE INDEX idx_cart_items_deleted_at ON cart_items(deleted_at);