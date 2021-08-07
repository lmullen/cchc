.PHONY : crawler
crawler : 
	docker build -f crawler.Dockerfile --tag crawler:test .

.PHONY : test-crawer
test-crawler : crawler
	docker run --rm \
		--env CCHC_DBHOST="host.docker.internal" \
		--env CCHC_LOGLEVEL="info" \
		-e CCHC_DBPORT \
		-e CCHC_DBUSER \
		-e CCHC_DBPASS \
		-e CCHC_DBNAME \
		--name crawler \
		crawler:test
