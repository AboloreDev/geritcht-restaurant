DROP INDEX IF EXISTS idx_orders_search_vector;

DROP TRIGGER IF EXISTS orders_search_vector_trigger ON orders;

DROP FUNCTION IF EXISTS update_orders_search_vector();

ALTER TABLE orders DROP COLUMN IF EXISTS search_vector;