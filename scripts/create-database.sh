#!/usr/bin/env bash -x

createdb --echo \
	--host=$CCHC_DBHOST \
	--port=$CCHC_DBPORT \
	--username=$CCHC_DBUSER \
	--owner=$CCHC_DBUSER \
	$CCHC_DBNAME
