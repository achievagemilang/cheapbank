postgres: 
	echo "Starting postgres"
	docker run --name cheapbank-postgres -p 54320:5432 -e POSTGRES_PASSWORD=postgres -d postgres 

createdb:
	echo "Creating database"
	docker exec -it cheapbank-postgres createdb --username=postgres --owner=postgres cheapbank

dropdb:
	echo "Dropping database"
	docker exec -it cheapbank-postgres dropdb cheapbank -U postgres

migrateup:
	echo "Migrating up"
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:54320/cheapbank?sslmode=disable" -verbose up

migratedown:
	echo "Migrating down"
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:54320/cheapbank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc