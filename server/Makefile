postgresinit:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root whisper-warp-db

dropdb:
	docker exec -it postgres15 dropdb whisper-warp-db

migrateup:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5433/whisper-warp-db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5433/whisper-warp-db?sslmode=disable" -verbose down

.PHONY: postgresinit postres createdb dropdb migrateup migratedown