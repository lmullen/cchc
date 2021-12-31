-- Create a table to hold the results of language detection
CREATE TABLE IF NOT EXISTS results.languages (
  job_id uuid REFERENCES jobs.fulltext (id) NOT NULL,
  item_id text REFERENCES items (id) NOT NULL,
  lang text NOT NULL,
  sentences integer
);

