include config.env
export

psql:
	sudo -u postgres psql $(POSTGRES_DB)

docker-psql:
	docker exec -it auditoria_bot_postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

docker-pgdump:
	docker exec -t auditoria_bot_postgres pg_dump --clean -U $(POSTGRES_USER) $(POSTGRES_DB) > backup/postgres/dump_`date +"%F_%T"`.sql

docker-pgrestore:
	cat $(DUMP_FILE) | docker exec -i auditoria_bot_postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)