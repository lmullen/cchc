# APPLICATION
# --------------------------------------------------

# Rebuild and run all services detached
.PHONY : up
up : 
	docker compose up --build --force-recreate --detach

# Run in production from containers
.PHONY : run
run :
	docker compose pull
	docker compose up --force-recreate --detach

.PHONY : stop
stop :
	docker compose stop

.PHONY : down
down :
	docker compose down

# DATABASE 
# --------------------------------------------------
.PHONY : migration,  db-up, db-down, db-drop

migration :
	@read -p "What is the slug for the migration? " migration;\
	migrate create -dir db/migrations -ext sql -seq $$migration

db-up :
	@echo "Migrating to current version of database"
	migrate -database "$(CCHC_DBSTR_LOCAL)" -path db/migrations up

db-down :
	migrate -database "$(CCHC_DBSTR_LOCAL)" -path db/migrations down 1

db-drop :
	@echo "Dropping the local database"
	migrate -database "$(CCHC_DBSTR_LOCAL)" -path migrations drop
