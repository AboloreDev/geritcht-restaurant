ALTER TABLE orders
ADD COLUMN search_vector tsvector;

CREATE INDEX idx_orders_search_vector
ON orders
USING GIN (search_vector);

CREATE OR REPLACE FUNCTION update_orders_search_vector()
RETURNS trigger AS $$
DECLARE
    u users%ROWTYPE;
BEGIN
    SELECT *
    INTO u
    FROM users
    WHERE id = NEW.user_id;

    NEW.search_vector :=
          setweight(
              to_tsvector('english',
                  coalesce(u.first_name, '') || ' ' ||
                  coalesce(u.last_name, '') || ' ' ||
                  coalesce(u.email, '')
              ),
              'A'
          )
        || setweight(
              to_tsvector('english',
                  coalesce(NEW.notes, '')
              ),
              'B'
          )
        || setweight(
              to_tsvector('english',
                  coalesce(NEW.status::text, '')
              ),
              'C'
          )
        || setweight(
              to_tsvector('english',
                  coalesce(NEW.type::text, '')
              ),
              'D'
          );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER orders_search_vector_trigger
BEFORE INSERT OR UPDATE
ON orders
FOR EACH ROW
EXECUTE FUNCTION update_orders_search_vector();

UPDATE orders
SET notes = notes;

COMMENT ON COLUMN orders.search_vector IS
'Full-text search vector.
A = User First Name, Last Name, Email
B = Search Notes
C = Order Status
D = Order Type';
