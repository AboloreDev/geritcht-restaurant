CREATE TABLE menu_categories (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL UNIQUE,
    description   TEXT,
    image_url     VARCHAR(500),
    display_order INT NOT NULL DEFAULT 0,
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_menu_categories_deleted_at ON menu_categories(deleted_at);
CREATE INDEX idx_menu_categories_is_active ON menu_categories(is_active);
CREATE INDEX idx_menu_categories_name ON menu_categories(name);
CREATE INDEX idx_menu_categories_created_at ON menu_categories(created_at);
CREATE INDEX idx_menu_categories_id ON menu_categories(id);
CREATE INDEX idx_menu_categories_id_is_active ON menu_categories(id, is_active) WHERE deleted_at IS NULL;