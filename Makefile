.PHONY : collection-crawler
collection-crawler : 
	docker build -f collection-crawler.Dockerfile --tag collection-crawler:test .

.PHONY : test-collection-crawer
test-collection-crawler : collection-crawler
	docker run --rm \
		--env CCHC_DBHOST="host.docker.internal" \
		-e CCHC_DBPORT \
		-e CCHC_DBUSER \
		-e CCHC_DBPASS \
		-e CCHC_DBNAME \
		--name collection-crawler \
		collection-crawler:test
