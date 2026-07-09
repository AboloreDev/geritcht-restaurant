DROP INDEX IF EXISTS idx_menus_search_vector;

DROP TRIGGER IF EXISTS menus_search_vector_trigger ON menus;

DROP FUNCTION IF EXISTS menu_search_vector_update();

ALTER TABLE menus DROP COLUMN IF EXISTS search_vector;