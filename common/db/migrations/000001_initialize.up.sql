-- Initialize the databse
CREATE TABLE IF NOT EXISTS collections (
  id text PRIMARY KEY,
  title text,
  description text,
  count integer,
  url text,
  items_url text,
  subjects text[],
  subjects2 text[],
  topics text[],
  api jsonb
);

CREATE TABLE IF NOT EXISTS items (
  id text PRIMARY KEY,
  url text,
  title text,
  year int,
  date text,
  subjects text[],
  fulltext text,
  fulltext_service text,
  fulltext_file text,
  timestamp bigint,
  api jsonb
);

CREATE TABLE IF NOT EXISTS items_in_collections (
  item_id text REFERENCES items (id),
  collection_id text REFERENCES collections (id),
  PRIMARY KEY (item_id, collection_id)
);

CREATE INDEX ON items (year);

CREATE INDEX ON items (timestamp);

CREATE INDEX ON items (subjects);

CREATE INDEX ON items (api)
WHERE
  api IS NULL;

