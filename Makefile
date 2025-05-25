include config.env
export

psql:
	sudo -u postgres psql $(POSTGRES_DB)

docker-psql:
	docker exec -it auditoria_bot_postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)
