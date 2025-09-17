
.PHONY: run test docker-build docker-up docker-test docker-down

run:
	go run main.go

test:
	go test ./... -v

docker-build:
	docker build -t config-management-service .

docker-up: docker-build
	docker run --name config-management-service -p 8090:8090 config-management-service

docker-down:
	docker stop config-management-service || true
	docker rm config-management-service || true

docker-test:
	ifeq ($(OS),Windows_NT)
		docker run --rm -v %cd%:/app -w /app golang:1.25 go test ./... -v
	else
		docker run --rm -v ${PWD}:/app -w /app golang:1.25 go test ./... -v
	endif
