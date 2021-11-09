ALTER TABLE jobs.fulltext_predict RENAME TO fulltext;

ALTER TABLE jobs.fulltext
  ADD COLUMN destination text;

UPDATE
  jobs.fulltext
SET
  destination = 'bible-quotations'
WHERE
  destination IS NULL;

ALTER TABLE jobs.fulltext
  ALTER COLUMN destination SET NOT NULL;

