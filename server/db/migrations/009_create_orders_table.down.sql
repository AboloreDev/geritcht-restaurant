ALTER TABLE order_items DROP CONSTRAINT IF EXISTS fk_order_items_order_id;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS order_type;