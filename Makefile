MIGRATION_FOLDER=$(CURDIR)/internal/migrations
POSTGRES_SETUP := user=postgres password=qwerty dbname=postgres host=localhost port=5432 sslmode=disable

build:
	docker compose build pvz

run:
	docker compose up pvz

migrate-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

migrate-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down

test-unit:
	go test ./... -tags=unit

test-integration:
	go test ./... -tags=integration

gen: 
	protoc --go_out=internal/pb/ \
		--go_opt=paths=import \
		--go-grpc_out=internal/pb/ \
		--go-grpc_opt=paths=import \
		api/order/order.proto \
		api/user/user.proto \
		api/pvz/pvz.proto \
		api/order_info/order_info.proto