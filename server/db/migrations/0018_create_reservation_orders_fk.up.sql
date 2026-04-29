-- add reservation foreign key to orders
-- done last because both tables needed to exist first
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_reservation_id
    FOREIGN KEY (reservation_id) REFERENCES reservations(id) ON DELETE SET NULL;