DROP INDEX IF EXISTS idx_reservations_search_vector;

DROP TRIGGER IF EXISTS reservations_search_vector_trigger ON reservations;

DROP FUNCTION IF EXISTS update_reservations_search_vector();

ALTER TABLE reservations DROP COLUMN IF EXISTS search_vector;