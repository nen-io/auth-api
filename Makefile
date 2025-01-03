
run:
	go run ./src/main.go

build:
	rm -rf ./build/* && go build -o build/api ./src/main.go

dev:
	air

db:
	docker compose up -d && go run ./src/cmd/migrate_db.go

clear_db:
	rm -rf db_data/*
