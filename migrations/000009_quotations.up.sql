CREATE SCHEMA IF NOT EXISTS results;

CREATE TABLE IF NOT EXISTS results.biblical_quotations (
  job_id uuid REFERENCES jobs.fulltext_predict (id) NOT NULL,
  item_id text REFERENCES items (id) NOT NULL,
  reference_id text NOT NULL,
  verse_id text NOT NULL,
  probability real NOT NULL
);

