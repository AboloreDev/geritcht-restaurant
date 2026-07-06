CREATE TABLE menus (
    id               BIGSERIAL PRIMARY KEY,
    menu_category_id      BIGINT NOT NULL REFERENCES menu_categories(id) ON DELETE RESTRICT,
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

CREATE INDEX idx_menus_menu_category_id ON menus(menu_category_id);
CREATE INDEX idx_menus_is_available ON menus(is_available);
CREATE INDEX idx_menus_deleted_at ON menus(deleted_at);
CREATE INDEX idx_menus_name_menu_category_id ON menus(name, menu_category_id);
CREATE INDEX idx_menus_id ON menus(id);
CREATE INDEX idx_menus_id_is_available ON menus(id, is_available) WHERE deleted_at IS NULL;
CREATE INDEX idx_menus_id_deleted_at ON menus(id, deleted_at);
CREATE INDEX idx_menus_name ON menus(name);
