CREATE TABLE menu_images (
    id           BIGSERIAL PRIMARY KEY,
    menu_item_id BIGINT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    url          VARCHAR(500) NOT NULL,
    alt_text     VARCHAR(255),
    is_primary   BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_menu_item_images_menu_item_id ON menu_images(menu_item_id);
CREATE INDEX idx_menu_item_images_deleted_at ON menu_images(deleted_at);