ALTER TABLE ingredients
ADD COLUMN search_vector TSVECTOR;

CREATE INDEX idx_ingredients_search_vector
ON ingredients
USING GIN(search_vector);

CREATE OR REPLACE FUNCTION ingredients_search_vector_update()
RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ingredients_search_vector_trigger
BEFORE INSERT OR UPDATE OF name
ON ingredients
FOR EACH ROW
EXECUTE FUNCTION ingredients_search_vector_update();

-- Populate existing rows by firing the trigger
UPDATE ingredients
SET name = name;

COMMENT ON COLUMN ingredients.search_vector IS
'Full-text search vector.
A = Ingredient Name';