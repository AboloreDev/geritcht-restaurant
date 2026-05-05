CREATE TABLE allergens (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE dietary_tags (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE menu_item_allergens (
    menu_id BIGINT NOT NULL REFERENCES menu(id) ON DELETE CASCADE,
    allergen_id  BIGINT NOT NULL REFERENCES allergens(id) ON DELETE CASCADE,
    PRIMARY KEY (menu_id, allergen_id)
);

CREATE TABLE menu_item_dietary (
    menu_id   BIGINT NOT NULL REFERENCES menu(id) ON DELETE CASCADE,
    dietary_tag_id BIGINT NOT NULL REFERENCES dietary_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (menu_id, dietary_tag_id)
);

CREATE INDEX idx_allergens_deleted_at ON allergens(deleted_at);
CREATE INDEX idx_dietary_tags_deleted_at ON dietary_tags(deleted_at);