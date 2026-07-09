ALTER TABLE menus ADD COLUMN search_vector TSVECTOR;

CREATE INDEX idx_menus_search_vector ON menus USING GIN(search_vector);

CREATE OR REPLACE FUNCTION menu_search_vector_update() RETURNS trigger AS $$
    BEGIN
        NEW.search_vector :=
            setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
            setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B');

        RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER menus_search_vector_trigger 
    BEFORE INSERT OR UPDATE OF name, description
    ON menus
    FOR EACH ROW
EXECUTE FUNCTION menu_search_vector_update();

-- Populate existing rows
UPDATE menus SET name = name;

COMMENT ON COLUMN menus.search_vector IS
'Full-text search vector.
A = Name
B = Description';