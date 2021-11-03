# APPLICATION
# --------------------------------------------------

# Rebuild and run all services detached
.PHONY : up
up : 
	docker compose up --build --force-recreate --detach

# Rebuild and run attached
.PHONY : attached
attached : 
	docker compose up --build --force-recreate

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

.PHONY: restart-qftext
restart-qftext :
	@echo "Restarting the full text enqueuer"
	docker compose stop qftext
	@mkdir -p logs
	docker compose logs qftext > logs/qftext-$(shell date +%FT%T).log
	docker compose up --build --detach qftext

.PHONY: restart-predictor
restart-predictor :
	@echo "Restarting the predictor"
	docker compose stop predictor
	@mkdir -p logs
	docker compose logs predictor > logs/predictor-$(shell date +%FT%T).log
	docker compose up --build --detach predictor

.PHONY: restart
restart : restart-itemmd restart-crawler restart-qftext restart-predictor

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
	docker compose logs qftext > logs/qftext-$(shell date +%FT%T).log
	docker compose logs predictor > logs/predictor-$(shell date +%FT%T).log
	docker compose down

# DATABASE
# --------------------------------------------------
.PHONY : db-create, db-up, db-down

db-up :
	migrate -database $(CCHC_DBSTR) -path db/migrations up

db-down :
	migrate -database $(CCHC_DBSTR) -path db/migrations down
