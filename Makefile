.PHONY : run

# APPLICATION
# --------------------------------------------------

# Rebuild and run all services detached
.PHONY : up
up : 
	docker compose up --build --force-recreate --detach

restart :
	@echo "Restarting just the crawler"
	docker compose stop crawler
	@mkdir -p logs
	docker compose logs crawer > logs/crawler-$(date +%FT%T).log
	docker compose up --build --detach crawler

.PHONY : stop
stop :
	docker compose stop

.PHONY : down
down :
	docker compose down

.PHONY : debug
debug :
	docker compose logs -f crawler

.PHONY : logs
logs :
	docker compose logs -f crawler | grep -v "level=debug"

.PHONY : collection-logs
collection-logs :
	docker compose logs -f crawler | grep "Fetched page of items from collection"

.PHONY : crawler
crawler : 
	docker compose build crawler

# DATABASE
# --------------------------------------------------
DBCONN="postgres://$(CCHC_DBUSER):$(CCHC_DBPASS)@$(CCHC_DBHOST):$(CCHC_DBPORT)/$(CCHC_DBNAME)?sslmode=disable"
.PHONY : db-create, db-up, db-down

db-create:
	./scripts/create-database.sh

db-up :
	migrate -database $(DBCONN) -path migrations up

db-down :
	migrate -database $(DBCONN) -path migrations down
