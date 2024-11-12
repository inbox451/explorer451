APP_NAME=explorer451

.PHONY: build run

build:
	go build -o bin/$(APP_NAME) main.go

run:
	go run --race main.go

test:
	curl http://localhost:8080/api/buckets/
