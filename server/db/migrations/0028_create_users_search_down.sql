DROP TABLE IF EXISTS idx_users_search_vector

DROP TRIGGER IF EXISTS users_search_vector_trigger ON users

DROP FUNCTION IF EXISTS users_search_vector_update()

ALTER TABLE users DROP COLUMN IF EXISTS search_vector