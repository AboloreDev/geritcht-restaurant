ALTER TABLE reservations
ADD COLUMN search_vector tsvector;

CREATE INDEX idx_reservations_search_vector
ON reservations
USING GIN (search_vector);

CREATE OR REPLACE FUNCTION update_reservations_search_vector()
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
                  coalesce(NEW.special_requests, '')
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
                  coalesce(NEW.date::text, '')
              ),
              'D'
          );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER reservations_search_vector_trigger
BEFORE INSERT OR UPDATE
ON reservations
FOR EACH ROW
EXECUTE FUNCTION update_reservations_search_vector();

UPDATE reservations
SET special_requests = special_requests;

COMMENT ON COLUMN reservations.search_vector IS
'Full-text search vector.
A = User First Name, Last Name, Email
B = Special Requests
C = Reservation Status
D = Reservation Date';
