CREATE VIEW table_sizes AS
SELECT
  n.nspname::text AS namespace,
  c.relname::text AS tablename,
  c.relkind AS type,
  pg_size_pretty(pg_total_relation_size(c.oid::regclass)) AS total_size,
  pg_total_relation_size(c.oid::regclass) AS size
FROM
  pg_class c
  LEFT JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE (n.nspname <> ALL (ARRAY['pg_catalog'::name, 'information_schema'::name]))
  AND c.relkind <> 'i'::"char"
  AND n.nspname !~ '^pg_toast'::text
ORDER BY
  (pg_total_relation_size(c.oid::regclass)) DESC;

ALTER VIEW table_sizes SET SCHEMA stats;

