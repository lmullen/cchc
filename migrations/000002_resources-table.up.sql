-- Move to a separate table for full text URLs
-- Drop the columns that used to store the full text
ALTER TABLE items
  DROP COLUMN IF EXISTS fulltext;

ALTER TABLE items
  DROP COLUMN IF EXISTS fulltext_service;

ALTER TABLE items
  DROP COLUMN IF EXISTS fulltext_file;

CREATE TABLE IF NOT EXISTS resources (
  item_id text REFERENCES items (id),
  resource_seq int,
  fulltext_file text,
  djvu_text_file text,
  image text,
  pdf text,
  url text,
  caption text,
  PRIMARY KEY (item_id, resource_seq)
);

CREATE TABLE IF NOT EXISTS files (
  item_id text REFERENCES items (id),
  resource_seq int,
  file_seq int,
  format_seq int,
  mimetype text,
  fulltext text,
  fulltext_service text,
  word_coordinates text,
  url text,
  info text,
  use text,
  PRIMARY KEY (item_id, resource_seq, file_seq, format_seq),
  FOREIGN KEY (item_id, resource_seq) REFERENCES resources (item_id, resource_seq)
);

CREATE INDEX ON resources (item_id);

CREATE INDEX ON files (item_id);

CREATE INDEX ON files (item_id, resource_seq);

