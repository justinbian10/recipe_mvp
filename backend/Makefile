include .envrc

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${RECIPE_MVP_DB_DSN} -cors-trusted-origins="http://localhost:5173"


.PHONY: db/psql
db/psql:
	psql ${RECIPE_MVP_DB_DSN}
