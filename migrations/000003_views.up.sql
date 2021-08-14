CREATE VIEW count_mimetypes AS
SELECT
  mimetype,
  count(*) AS n
FROM
  files
GROUP BY
  mimetype
ORDER BY
  n DESC;

