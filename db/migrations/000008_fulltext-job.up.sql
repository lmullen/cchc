CREATE SCHEMA IF NOT EXISTS jobs;

CREATE TYPE text_level AS ENUM (
  'file',
  'resource'
);

CREATE TABLE IF NOT EXISTS jobs.fulltext_predict (
  id uuid PRIMARY KEY,
  item_id text REFERENCES items (id) NOT NULL,
  level text_level,
  source text,
  has_ft_method boolean NOT NULL,
  started timestamp,
  finished timestamp
);

CREATE INDEX ON jobs.fulltext_predict (item_id);

CREATE VIEW jobs.fulltext_unqueued AS
SELECT
  i.id
FROM
  items i
WHERE
  i.api IS NOT NULL
  AND NOT EXISTS (
    SELECT
    FROM
      jobs.fulltext_predict
    WHERE
      fulltext_predict.item_id = i.id);

CREATE VIEW jobs.fulltext_job_1_skipped AS
SELECT
  *
FROM
  jobs.fulltext_predict
WHERE
  NOT has_ft_method;

CREATE VIEW jobs.fulltext_job_2_started AS
SELECT
  *
FROM
  jobs.fulltext_predict
WHERE
  started IS NOT NULL
  AND finished IS NULL;

CREATE VIEW jobs.fulltext_job_3_finished AS
SELECT
  *
FROM
  jobs.fulltext_predict
WHERE
  started IS NOT NULL
  AND finished IS NOT NULL;

