DROP INDEX IF EXISTS items_subjects_idx;

CREATE INDEX IF NOT EXISTS items_subjects_idx ON items USING btree (subjects);

DROP INDEX IF EXISTS files_mimetype_idx;

