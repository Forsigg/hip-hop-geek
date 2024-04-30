test:
	go test ./internal/...

run:
	go run cmd/releases/main.go

coverage:
	go test -cover -coverprofile coverage.out ./internal/...

db_status:
	cd ./internal/db/migrations/ && goose sqlite ../db.db status && cd ../../

db_up:
	cd ./internal/db/migrations/ && goose sqlite ../db.db up && cd ../../

db_down:
	cd ./internal/db/migrations/ && goose sqlite ../db.db down && cd ../../
