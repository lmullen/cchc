-- Create job status better
ALTER TABLE jobs.fulltext
  DROP COLUMN IF EXISTS has_ft_method;

CREATE TYPE job_statuses AS ENUM (
  'ready',
  'running',
  'skipped',
  'failed',
  'finished'
);

ALTER TABLE jobs.fulltext
  ADD COLUMN IF NOT EXISTS status job_statuses;

CREATE INDEX IF NOT EXISTS jobs_status_idx ON jobs.fulltext USING btree (status);

