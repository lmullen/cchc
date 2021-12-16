DROP INDEX IF EXISTS idx_items_languages;

ALTER TABLE items DROP COLUMN IF EXISTS languages;

