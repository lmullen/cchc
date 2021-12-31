ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS resource_seq integer; 

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS file_seq integer; 

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS format_seq integer; 

CREATE TYPE text_level AS ENUM (
  'file',
  'resource'
);

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS format_seq text_level; 

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS source text; 

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS source text; 

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS destination;

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS queue text; 

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS started;

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS started timestamp without time zone;

ALTER TABLE jobs.fulltext DROP COLUMN IF EXISTS finished;

ALTER TABLE jobs.fulltext ADD COLUMN IF NOT EXISTS finished timestamp without time zone;
