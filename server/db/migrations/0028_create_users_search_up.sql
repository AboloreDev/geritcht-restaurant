ALTER TABLE users
ADD COLUMN search_vector TSVECTOR;

CREATE INDEX idx_users_search_vector
ON users
USING GIN(search_vector);

CREATE OR REPLACE FUNCTION users_search_vector_update()
RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.email, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.first_name, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(NEW.last_name, '')), 'C');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_search_vector_trigger
BEFORE INSERT OR UPDATE OF email, first_name, last_name
ON users
FOR EACH ROW
EXECUTE FUNCTION users_search_vector_update();

-- Populate existing rows
UPDATE users
SET email = email;

CREATE OR REPLACE FUNCTION refresh_user_orders_search_vector()
RETURNS trigger AS $$
BEGIN
    UPDATE orders
    SET search_vector =
          setweight(
              to_tsvector('english',
                  coalesce(NEW.first_name, '') || ' ' ||
                  coalesce(NEW.last_name, '') || ' ' ||
                  coalesce(NEW.email, '')
              ),
              'A'
          )
        || setweight(
              to_tsvector('english',
                  coalesce(notes, '')
              ),
              'B'
          )
        || setweight(
              to_tsvector('english',
                  coalesce(status::text, '')
              ),
              'C'
          )
        || setweight(
              to_tsvector('english',
                  coalesce(type::text, '')
              ),
              'D'
          )
    WHERE user_id = NEW.id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_search_vector_trigger
AFTER UPDATE OF first_name, last_name, email
ON users
FOR EACH ROW
EXECUTE FUNCTION refresh_user_orders_search_vector();

COMMENT ON COLUMN users.search_vector IS
'Full-text user search vector.
A = Email
B = First Name
C = Last Name';