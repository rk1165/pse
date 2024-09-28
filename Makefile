run:
	go run --tags "fts5" ./cmd/web/*.go

init:
	sqlite3 engine.db < ddl.sql
	go mod tidy

build:
	go build --tags "fts5" -o build/pse ./cmd/web/*.go

.PHONY: run init build