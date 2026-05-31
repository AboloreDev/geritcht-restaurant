CREATE UNIQUE INDEX idx_reservations_unique_booking
ON reservations(table_id, date, time_slot)
WHERE status NOT IN ('cancelled', 'no_show');
