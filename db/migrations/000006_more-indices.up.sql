CREATE INDEX IF NOT EXISTS files_mimetype_idx ON files USING btree (mimetype);

DROP INDEX IF EXISTS items_subjects_idx;

CREATE INDEX IF NOT EXISTS items_subjects_idx ON items USING gin (subjects);

