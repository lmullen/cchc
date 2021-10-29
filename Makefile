# APPLICATION
# --------------------------------------------------

# Rebuild and run all services detached
.PHONY : up
up : 
	docker compose up --build --force-recreate --detach

.PHONY: restart-crawler
restart-crawler :
	@echo "Restarting the crawler"
	docker compose stop crawler
	@mkdir -p logs
	docker compose logs crawler > logs/crawler-$(shell date +%FT%T).log
	docker compose up --build --detach crawler

.PHONY: restart-itemmd
restart-itemmd :
	@echo "Restarting the item metadata fetcher"
	docker compose stop itemmd
	@mkdir -p logs
	docker compose logs itemmd > logs/itemmd-$(shell date +%FT%T).log
	docker compose up --build --detach itemmd

.PHONY: restart
restart : restart-itemmd restart-crawler

.PHONY : stop
stop :
	docker compose stop

.PHONY : down
down :
	docker compose stop
	@mkdir -p logs
	docker compose logs crawler > logs/crawler-$(shell date +%FT%T).log
	docker compose logs itemmd > logs/itemmd-$(shell date +%FT%T).log
	docker compose logs queue > logs/queue-$(shell date +%FT%T).log
	docker compose logs predictor > logs/predictor-$(shell date +%FT%T).log
	docker compose down

# DATABASE
# --------------------------------------------------
DBHOST=$(CCHC_DBHOST)
# If using Docker on a Mac, set DBHOST to localhost
ifeq ($(DBHOST), host.docker.internal)
DBHOST=localhost
endif
DBCONN="postgres://$(CCHC_DBUSER):$(CCHC_DBPASS)@$(DBHOST):$(CCHC_DBPORT)/$(CCHC_DBNAME)?sslmode=disable"
.PHONY : db-create, db-up, db-down

db-create:
	./scripts/create-database.sh

db-up :
	migrate -database $(DBCONN) -path db/migrations up

db-down :
	migrate -database $(DBCONN) -path db/migrations down
