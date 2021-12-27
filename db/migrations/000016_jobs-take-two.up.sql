-- Replace the old version of the jobs table with a simpler one
DROP VIEW IF EXISTS jobs.fulltext_job_1_skipped;

DROP VIEW IF EXISTS jobs.fulltext_job_2_started;

DROP VIEW IF EXISTS jobs.fulltext_job_3_finished;

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS resource_seq; 

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS file_seq; 

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS format_seq; 

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS level; 

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS source; 

DROP TYPE IF EXISTS text_level;

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS queue;

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS destination text;

ALTER TABLE jobs.fulltext ALTER COLUMN destination SET NOT NULL;

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS started;

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS started timestamp with time zone;

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS finished;

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS finished timestamp with time zone;
