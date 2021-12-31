-- Add a view with job status
CREATE VIEW job_status_ft AS
SELECT
  destination,
  status,
  COUNT(*) AS num_jobs
FROM
  jobs.fulltext
GROUP BY
  destination,
  status
ORDER BY
  destination,
  status;

ALTER VIEW job_status_ft SET SCHEMA stats;

