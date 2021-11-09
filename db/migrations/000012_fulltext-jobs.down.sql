ALTER TABLE jobs.fulltext
  DROP COLUMN destination;

ALTER TABLE jobs.fulltext RENAME TO fulltext_predict;

