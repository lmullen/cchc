CREATE SCHEMA IF NOT EXISTS stats;

ALTER VIEW count_mimetypes SET SCHEMA stats;

CREATE VIEW items_per_collection AS
SELECT
  collections.title,
  COUNT(items_in_collections.collection_id) AS n_items,
  collections.id
FROM
  collections
  LEFT JOIN items_in_collections ON collections.id = items_in_collections.collection_id
GROUP BY
  collections.id
ORDER BY
  n_items DESC;

ALTER VIEW items_per_collection SET SCHEMA stats;

CREATE VIEW items_status AS (
  SELECT
    'total items' AS type,
    COUNT(*) AS n
  FROM
    items)
UNION (
  SELECT
    'fetched items' AS type,
    COUNT(*) AS n
  FROM
    items
  WHERE
    api IS NOT NULL)
UNION (
  SELECT
    'unfetched items' AS type,
    COUNT(*) AS n
  FROM
    items
  WHERE
    api IS NULL)
ORDER BY
  n DESC;

ALTER VIEW items_status SET SCHEMA stats;

