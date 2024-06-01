postgres:
	docker run --name postgres-alpine -p 5432:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:alpine

createdb:
	docker exec -it postgres-alpine createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-alpine dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...	

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple_bank/db/sqlc IStore	

.PHONY: createdb dropdb postgres migratedown migrateup migratedown1 migrateup1 sqlc test server mock
