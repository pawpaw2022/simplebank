postgres: 
	docker run --name postgres12 --network bank-network -p:5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb: 
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank 

dropdb:
	docker exec -it postgres12 dropdb simple_bank

# migrateup:
# 	 migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup:
	 migrate -path db/migrations -database "postgresql://root:4qrq04csWVrXTMXXlSIq@simple-bank.ccehaovvsjiq.us-east-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migratedown:
	 migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	 migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	 migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/pawpaw2022/simplebank/db/postgresql Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock