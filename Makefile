run:
	go run --tags "fts5" ./cmd/web/

init:
	sqlite3 engine.db < ./sql/ddl.sql
	go mod tidy

build:
	go build --tags "fts5" -o pse ./cmd/web/

clean:
	sqlite3 engine.db < ./sql/clean.sql

.PHONY: run init build clean