DROP INDEX IF EXISTS jobs_status_idx;

ALTER TABLE jobs.fulltext
  DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS job_statuses;

ALTER TABLE jobs.fulltext
  ADD COLUMN IF NOT EXISTS has_ft_method bool;

