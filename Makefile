# APPLICATION
# --------------------------------------------------

# Rebuild and run all services detached
.PHONY : build, up, down
build :
	docker compose --profile db --profile api --profile languages --profile quotations build

up : 
	docker compose --profile api --profile languages up --build --force-recreate --detach

# Stops ALL profiles
down :
	docker compose --profile db --profile api --profile languages --profile quotations down

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


# DEPLOY 
# --------------------------------------------------
.PHONY : deploy

deploy : export CCHC_VERSION=release
deploy : 
	docker compose --profile ctrl build --parallel
	docker push ghcr.io/lmullen/cchc-crawler:release
	docker push ghcr.io/lmullen/cchc-itemmd:release
	docker push ghcr.io/lmullen/cchc-ctrl:release
	docker push ghcr.io/lmullen/cchc-language-detector:release
	docker push ghcr.io/lmullen/cchc-predictor:release
