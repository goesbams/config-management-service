
.PHONY: run test docker-build docker-up docker-down

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
