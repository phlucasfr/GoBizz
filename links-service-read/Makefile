createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres gobizz

dropdb:
	docker exec -it postgres dropdb gobizz

migrateup:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/gobizz?sslmode=disable" -verbose up;      

migratedown:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/gobizz?sslmode=disable" -verbose down;

sqlc:
	sqlc generate

test:
	go test -v -cover ./... 

server:
	go run cmd/main.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/links_read.proto

.PHONY: createdb, dropdb, migrateup, migratedown, sqlc, test, server, proto