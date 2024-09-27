run:
	go run --tags "fts5" ./cmd/web/*.go

init:
	sqlite3 engine.db < ddl.sql
	go mod tidy

.PHONY: run init