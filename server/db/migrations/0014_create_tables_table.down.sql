ALTER TABLE reservations DROP CONSTRAINT IF EXISTS fk_reservations_table_id;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_table_id;
DROP TABLE IF EXISTS tables;
DROP TYPE IF EXISTS table_status;