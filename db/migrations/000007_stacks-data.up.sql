CREATE TABLE IF NOT EXISTS stacks_books (
  lccn text PRIMARY KEY,
  isbn text[],
  title text,
  publisher text,
  date date,
  year int,
  subject_full text[],
  subject text[],
  person text[],
  lang text[],
  original_metadata jsonb,
  text text
);

CREATE INDEX IF NOT EXISTS stacks_books_year_idx ON stacks_books USING btree (year);

CREATE INDEX IF NOT EXISTS stacks_books_subject_idx ON stacks_books USING gin (subject);

CREATE INDEX IF NOT EXISTS stacks_books_subject_full_idx ON stacks_books USING gin (subject_full);

