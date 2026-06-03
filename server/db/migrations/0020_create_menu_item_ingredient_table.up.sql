CREATE TABLE menu_item_ingredients (
    menu_id  BIGINT NOT NULL REFERENCES menu(id) ON DELETE CASCADE,
    ingredient_id BIGINT NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity      DECIMAL(10,2) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (menu_id, ingredient_id)
);

CREATE INDEX idx_menu_item_ingredients_menu_id ON menu_item_ingredients(menu_id);
CREATE INDEX idx_menu_item_ingredients_ingredient_id ON menu_item_ingredients(ingredient_id);