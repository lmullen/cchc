CREATE INDEX IF NOT EXISTS quotations_job_id_idx ON results.biblical_quotations (job_id);

CREATE INDEX IF NOT EXISTS languages_job_id_idx ON results.languages (job_id);

CREATE INDEX IF NOT EXISTS quotations_item_id_idx ON results.biblical_quotations (item_id);

CREATE INDEX IF NOT EXISTS languages_item_id_idx ON results.languages (item_id);

