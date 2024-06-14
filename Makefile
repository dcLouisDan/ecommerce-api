build:
	@go build -o web_projects/ecommerce-api cmd/main.go

test:
	@go test -v ./...

run: build
	@./web_projects/ecommerce-api

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down
