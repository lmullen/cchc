# APPLICATION
# --------------------------------------------------

# Rebuild and run all services detached
.PHONY : up
up : 
	docker compose up --build --force-recreate --detach

# DATABASE 
# --------------------------------------------------
.PHONY : migration,  db-up, db-down

migration :
	@read -p "What is the slug for the migration? " migration;\
	migrate create -dir common/db/migrations -ext sql -seq $$migration

db-up :
	@echo "Migrating to current version of database"
	migrate -database "$(CCHC_DBSTR_LOCAL)" -path common/db/migrations up

db-down :
	migrate -database "$(CCHC_DBSTR_LOCAL)" -path common/db/migrations down 1
