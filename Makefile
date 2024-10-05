run:
	go run --tags "fts5" ./cmd/web/

init:
	sqlite3 engine.db < ddl.sql
	go mod tidy

build:
	go build --tags "fts5" -o pse ./cmd/web/

clean:
	sqlite3 engine.db < clean.sql

.PHONY: run init build clean