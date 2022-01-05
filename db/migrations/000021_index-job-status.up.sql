CREATE INDEX IF NOT EXISTS jobs_destination_idx ON jobs.fulltext (destination);

CREATE INDEX IF NOT EXISTS jobs_destination_status_idx ON jobs.fulltext (destination, status);

