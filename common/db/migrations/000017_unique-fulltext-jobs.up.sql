-- Make sure there is only one job for each combination of item and destination
ALTER TABLE jobs.fulltext
  ADD CONSTRAINT jobs_fulltext_unique UNIQUE (item_id, destination);

