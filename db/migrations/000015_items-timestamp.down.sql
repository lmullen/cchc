DROP INDEX IF EXISTS idx_items_updated;

ALTER TABLE items DROP COLUMN IF EXISTS updated;
