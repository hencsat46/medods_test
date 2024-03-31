doc:
	docker-compose up
build:
	@go build -o ./bin/main ./cmd/main.go
img:
	docker buildx build .