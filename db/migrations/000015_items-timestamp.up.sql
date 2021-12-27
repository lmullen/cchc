ALTER TABLE items ADD COLUMN IF NOT EXISTS updated TIMESTAMP WITH TIME ZONE;

-- Create data for items that already exist and have been fetched
UPDATE items SET updated = NOW();

ALTER TABLE items ALTER COLUMN updated SET NOT NULL;

CREATE INDEX idx_items_updated ON items USING brin(updated);
