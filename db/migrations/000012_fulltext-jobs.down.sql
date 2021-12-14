ALTER TABLE jobs.fulltext
  DROP COLUMN queue;

ALTER TABLE jobs.fulltext RENAME TO fulltext_predict;

