CREATE SCHEMA IF NOT EXISTS results;

CREATE TABLE IF NOT EXISTS results.biblical_quotations (
  job_id item_id,
  reference_id verse_id,
  probability
