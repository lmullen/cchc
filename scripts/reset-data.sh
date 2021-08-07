#!/usr/bin/env bash -x

# Drop the database tables with metadata
psql -e \
  -d "host=$CCHC_DBHOST port=$CCHC_DBPORT user=$CCHC_DBUSER dbname=$CCHC_DBNAME" \
	<< SQLSCRIPT
	DROP TABLE IF EXISTS collections CASCADE;
	DROP TABLE IF EXISTS items CASCADE;
	DROP TABLE IF EXISTS items_in_collections CASCADE;
SQLSCRIPT

# Drop the container with the message queue
docker rm -f cchc-queue

# Drop the volume which contains the persistent message queue
docker volume rm cchc_queue-data || true
