postgres:
	docker-compose up

createdb: 
	docker exec -it simple-bank createdb --username=yinnohs --owner=yinnohs simple_bank

dropdb:
	docker exec -it simple-bank dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://yinnohs:1234@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://yinnohs:1234@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down
	
sqlcgenerate:
	sqlc generate

test:
	go test -v -cover ./...

dep:
	go mod tidy

.PHONY: postgres createdb dropdb migrateup migratedown sqlcgenerate dep