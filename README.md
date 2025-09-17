# Config Management Service

A service to create, update, fetch, rollback, and manage versions of configurations.

---

## Table of Contents

- [Requirements](#requirements)  
- [Running Locally](#running-locally)  
  - [Using Go](#using-go)  
  - [Using Docker](#using-docker)  
- [API Endpoints](#api-endpoints)  
  - [Create a new config](#1-create-a-new-config)  
  - [Update a config](#2-update-a-config)  
  - [Rollback a config](#3-rollback-a-config)  
  - [Fetch a config](#4-fetch-a-config)  
  - [List all versions](#5-list-all-versions)  
- [OpenAPI Specification](#openapi-specification)  

---

## Requirements

- Go >= 1.25  
- Docker (optional for containerized setup)  
- curl (for testing APIs)
- `make` installed

---

## Running Locally

### Using Go

1. Clone the repository:

```sh
git clone https://github.com/goesbams/config-management-service.git
cd config-management-service
```

2. Run the service
```sh
make run
```
or
```sh
go run main.go
```

3. Run test
```
make test
```

### Using Docker

1. Build docker image:

```
make docker-build
```
2. Start container
```
make docker-up
```
3. Stop and remove container
```
make docker-down
```

## API Endpoints
### 1. Create a new config

```bash
curl -X POST http://localhost:8090/config \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d '{"name":"Main Database Config","type":"DATABASE","versions":[{"version":1,"property":{"max_limit":1000,"enabled":true}}]}'
```

### 2. Update a config

```bash
curl -X POST http://localhost:8090/config/update \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d '{"name":"Main Database Config","type":"DATABASE","versions":[{"version":2,"property":{"max_limit":2000,"enabled":false}}]}'
```

### 3. Rollback a config

```bash
curl -X POST http://localhost:8090/config/rollback \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d '{"name":"Main Database Config","version":1}'
```

### 4. Fetch a config

- Latest version:
```bash
curl -X GET "http://localhost:8090/config/fetch?name=Main%20Database%20Config" \
  -H "accept: application/json"
```

- Specific version:
```bash
curl -X GET "http://localhost:8090/config/fetch?name=Main%20Database%20Config&version=2" \
  -H "accept: application/json"
```

### 5. List all versions

```bash
curl -X GET "http://localhost:8090/config/versions?name=Main%20Database%20Config" \
  -H "accept: application/json"
```

## OpenAPI Specification
- All endpoints, request bodies, and responses are documented there.
- Use it with Swagger UI, code generators, or API clients.

