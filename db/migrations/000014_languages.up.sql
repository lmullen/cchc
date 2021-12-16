ALTER TABLE items
	ADD COLUMN IF NOT EXISTS languages text[];

-- Create data for items that already exist and have been fetched
UPDATE
	items
SET
	languages = ARRAY (
		SELECT
			jsonb_array_elements_text(api -> 'item' -> 'language'))
WHERE
	api IS NOT NULL;

CREATE INDEX idx_items_languages ON items USING GIN(languages);
