CREATE TABLE menu (
    id               BIGSERIAL PRIMARY KEY,
    category_id      BIGINT NOT NULL REFERENCES menu_categories(id) ON DELETE RESTRICT,
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    price            DECIMAL(10,2) NOT NULL,
    image_url        VARCHAR(500),
    is_available     BOOLEAN NOT NULL DEFAULT true,
    prep_time_minutes INT NOT NULL DEFAULT 15,
    spice_level      INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_menu_items_category_id ON menu(category_id);
CREATE INDEX idx_menu_items_is_available ON menu(is_available);
CREATE INDEX idx_menu_items_deleted_at ON menu(deleted_at);