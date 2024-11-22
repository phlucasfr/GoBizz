migrateup:
	migrate -path migrations -database "your-pg-conn-string" -verbose up;      

migratedown:
	migrate -path migrations -database "your-pg-conn-string" -verbose down;

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/main.go

.PHONY: migrateup, migratedown, sqlc, test, server
