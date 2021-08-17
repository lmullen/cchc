# APPLICATION
# --------------------------------------------------

# Rebuild and run all services detached
.PHONY : up
up : 
	docker compose up --build --force-recreate --detach

restart :
	@echo "Restarting the crawler and the item metadata fetcher"
	docker compose stop crawler
	docker compose stop itemmd
	@mkdir -p logs
	docker compose logs crawer > logs/crawler-$(shell date +%FT%T).log
	docker compose logs itemmd > logs/itemmd-$(shell date +%FT%T).log
	docker compose up --build --detach crawler
	docker compose up --build --detach itemmd

.PHONY : stop
stop :
	docker compose stop

.PHONY : down
down :
	docker compose stop
	docker compose logs crawler > logs/crawler-$(shell date +%FT%T).log
	docker compose logs itemmd > logs/itemmd-$(shell date +%FT%T).log
	docker compose logs queue > logs/queue-$(shell date +%FT%T).log
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
	migrate -database $(DBCONN) -path migrations up

db-down :
	migrate -database $(DBCONN) -path migrations down
