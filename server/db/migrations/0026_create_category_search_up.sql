ALTER TABLE menu_categories ADD COLUMN search_vector TSVECTOR;

CREATE INDEX idx_menu_categories_search_vector ON menu_categories USING GIN(search_vector);

CREATE OR REPLACE FUNCTION menu_categories_search_vector_update()
    RETURNS trigger AS $$
    BEGIN
        NEW.search_vector :=
            setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
            setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER menu_categories_search_vector_trigger BEFORE INSERT OR UPDATE OF name, description
    ON menu_categories
    FOR EACH ROW
EXECUTE FUNCTION menu_categories_search_vector_update();

-- Populate existing rows by firing the trigger
UPDATE menu_categories SET name = name;

COMMENT ON COLUMN menu_categories.search_vector IS
'Full-text menu category search vector.
A = Name
B = Description';