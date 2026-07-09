DROP INDEX IF EXISTS idx_ingredients_search_vector;

DROP TRIGGER IF EXISTS ingredients_search_vector_trigger ON ingredients;

DROP FUNCTION IF EXISTS ingredients_search_vector_update();

ALTER TABLE ingredients DROP COLUMN IF EXISTS search_vector;