-- Return the database to empty
BEGIN;
DROP TABLE IF EXISTS collections CASCADE;
DROP TABLE IF EXISTS items CASCADE;
DROP TABLE IF EXISTS items_in_collections CASCADE;
END;
