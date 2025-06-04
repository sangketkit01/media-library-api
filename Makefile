DB_URL=postgresql://root:secret@localhost:5433/media_library?sslmode=disable

createdb:
	docker exec -it postgres createdb --username=root --owner=root media_library

dropdb:
	docker exec -it postgres dropdb --username=root --owner=root media_library

migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir internal/db/migration $(name)

server:
	cd cmd/api && go run .

sqlc:
	sqlc generate

.proto:
	rm -f pb/*.go
	protoc \
	--proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	proto/*.proto

.PHONY: createdb dropdb migrateup migratedown new_migration server sqlc