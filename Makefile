.DEFAULT_GOAL := app
app:
	docker compose down && docker compose up --build
demo-client:
	go run cmd/client/main.go -input=input-B -output=output-B
lint:
	gofmt -w -s . && go mod tidy && go vet ./...
