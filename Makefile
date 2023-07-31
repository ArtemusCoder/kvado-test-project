.SILENT:


build:
	go mod download && docker-compose -f ./database/docker-compose.yml build

run_db:
	docker-compose -f ./database/docker-compose.yml up

stop_db:
	docker-compose -f ./database/docker-compose.yml down

run_server:
	go run ./cmd/server/main.go


clean_db:
	docker-compose -f ./database/docker-compose.yml down -v --rmi all --remove-orphans


tests:
	go test ./test -v