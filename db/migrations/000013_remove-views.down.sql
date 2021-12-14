-- Taken from 000008 up
CREATE VIEW jobs.fulltext_unqueued AS
SELECT
  i.id
FROM
  items i
WHERE
  i.api IS NOT NULL
  AND NOT EXISTS (
    SELECT
    FROM
      jobs.fulltext
    WHERE
      fulltext.item_id = i.id);

CREATE VIEW jobs.fulltext_queued AS
SELECT
  i.id
FROM
  items i
WHERE
  i.api IS NOT NULL
  AND (EXISTS (
      SELECT
      FROM
        jobs.fulltext
      WHERE
        fulltext.item_id = i.id));

