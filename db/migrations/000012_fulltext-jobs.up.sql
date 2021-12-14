ALTER TABLE jobs.fulltext_predict RENAME TO fulltext;

ALTER TABLE jobs.fulltext
  ADD COLUMN queue text;

