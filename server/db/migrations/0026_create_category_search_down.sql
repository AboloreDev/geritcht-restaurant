DROP INDEX IF EXISTS idx_menu_categories_search_vector;

DROP TRIGGER IF EXISTS menu_categories_search_vector_trigger ON menu_categories;

DROP FUNCTION IF EXISTS menu_categories_search_vector_update();

ALTER TABLE menu_categories DROP COLUMN IF EXISTS search_vector;
