postgresinit:
	docker run --name whisper-warp-postgres-local -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:15-alpine

postgres:
	docker exec -it whisper-warp-postgres-local psql -h localhost -p 5433 -U root -d whisper-warp-db

createdb:
	docker exec -it whisper-warp-postgres-local createdb --username=root --owner=root whisper-warp-db

dropdb:
	docker exec -it whisper-warp-postgres-local dropdb whisper-warp-db

migrateup:
	migrate -path db/new_migrations -database "postgresql://root:admin@localhost:5433/whisper-warp-db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/new_migrations -database "postgresql://root:admin@localhost:5433/whisper-warp-db?sslmode=disable" -verbose down

.PHONY: postgresinit postres createdb dropdb migrateup migratedown