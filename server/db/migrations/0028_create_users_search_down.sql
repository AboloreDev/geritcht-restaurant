DROP INDEX IF EXISTS idx_users_search_vector;

DROP TRIGGER IF EXISTS users_search_vector_trigger ON users;

DROP FUNCTION IF EXISTS users_search_vector_update();
DROP FUNCTION IF EXISTS refresh_user_orders_search_vector();

ALTER TABLE users DROP COLUMN IF EXISTS search_vector;