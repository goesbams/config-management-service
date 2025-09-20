.PHONY: help docker-build docker-up docker-down docker-openapi docker-stop docker-logs ps

SERVICE_APP := config-management-app
SERVICE_SWAGGER := swagger
OPENAPI := api/openapi.yaml
SERVICE ?= $(SERVICE_APP)

help:
	@echo "Targets:"
	@echo "  docker-build     Build $(SERVICE_APP) image"
	@echo "  docker-up        Build start ONLY $(SERVICE_APP)"
	@echo "  docker-down      Stop & remove ONLY $(SERVICE_APP)"
	@echo "  docker-openapi   Start ONLY $(SERVICE_SWAGGER)"
	@echo "  docker-stop      docker compose down (semua)"
	@echo "  docker-logs      Tail logs (default: $(SERVICE_APP))"
	@echo "  ps               Tampilkan status container"

docker-build:
	docker compose build $(SERVICE_APP)

docker-up: docker-build
	docker compose up --no-deps -d $(SERVICE_APP)

docker-down:
	-docker compose stop $(SERVICE_APP)
	-docker compose rm -f $(SERVICE_APP)

docker-openapi:
	test -f "$(OPENAPI)" || { echo "Missing $(OPENAPI)"; exit 1; }
	docker compose up --no-deps -d $(SERVICE_SWAGGER)

docker-stop:
	docker compose down

docker-logs:
	docker compose logs -f $(SERVICE)

ps:
	docker compose ps
