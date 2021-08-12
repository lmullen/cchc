.PHONY : run

# APPLICATION
# --------------------------------------------------

# Rebuild and run all services detached
run : 
	docker compose up -d --build --force-recreate --detach

.PHONY : logs
logs :
	docker compose logs -f crawler

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
