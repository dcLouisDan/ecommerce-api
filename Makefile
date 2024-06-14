build:
	@go build -o web_projects/ecommerce-api cmd/main.go

test:
	@go test -v ./...

run: build
	@./web_projects/ecommerce-api
