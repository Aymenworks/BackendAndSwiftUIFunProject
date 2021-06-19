help:
	@cat Makefile

run:
	@go run main.go

build:
	@go run main.go

db-create:
	@docker exec -i localstack-db mysql -e "\
	DROP DATABASE IF EXISTS ${COOKING_TIPS_MYSQL_DATABASE}; \
	CREATE DATABASE ${COOKING_TIPS_MYSQL_DATABASE} DEFAULT COLLATE utf8mb4_general_ci;"

# TODO: create the table and fields as well there
# TODO: think of automatic migration as well in the future with up/down mechanisms