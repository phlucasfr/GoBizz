migrateup:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up;      

migratedown:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose down;

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/main.go

.PHONY: migrateup, migratedown, sqlc, test, server
