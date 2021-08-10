.PHONY : run
# Rebuild and run all services detached
run : 
	docker compose up -d --build --force-recreate --detach

.PHONY : logs
logs :
	docker compose logs -f crawler

.PHONY : crawler
crawler : 
	docker compose build crawler

