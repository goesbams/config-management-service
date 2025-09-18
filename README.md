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
- [Schema Explanation](#schema-explanation)
  - [Config Object](#config-object)
  - [Versions Object](#version-object)
- [Design Decisions & Trade-offs](#design-decisions--trade-offs)
  - [Versioning per config](#1-versioning-per-config)
  - [In-memory vs persistent storage](#2-in-memory-vs-persistent-storage)
  - [Dockerized setup](#3-dockerized-setup)
- [Potential Improvements & Future Features](#potential-improvements--future-features)
- [OpenAPI Specification](#openapi-specification)
  - [How to access openapi-swagger](#how-to-access-openapi-swagger)

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
```sh
make test
```

### Using Docker

1. Build docker image:

```sh
make docker-build
```
2. Start container
```sh
make docker-up
```
3. Stop and remove container
```sh
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

## Schema Explanation
### Config Object
| Field      | Type   | Description                                        |
| ---------- | ------ | -------------------------------------------------- |
| `name`     | string | Unique name of the configuration                   |
| `type`     | string | Type of configuration (`DATABASE`, `API`, etc.)    |
| `versions` | array  | List of version objects (see Version Object below) |

### Versions Object
| Field      | Type   | Description                                 |
| ---------- | ------ | ------------------------------------------- |
| `version`  | int    | Version number                              |
| `property` | object | Key-value pairs of configuration properties |


## Design Decisions & Trade-offs

### 1. Versioning per config
- Allows rollback and history tracking.
- Trade-off: extra storage overhead for large configs.

### 2. In-memory vs persistent storage
- Currently in-memory for simplicity and test purposes.
- Persistent DB could be added for durability.

### 3. Dockerized setup
- Ensures environment consistency.
- Trade-off: additional Docker knowledge required.

## Potential Improvements & Future Features
- Add persistent storage with PostgreSQL
- Integration with CI/CD pipelines

## OpenAPI Specification
- All endpoints, request bodies, and responses are documented there.
- Use it with Swagger UI, code generators, or API clients.

### How to access openapi-swagger

1. Run the command
```bash
make docker-openapi
```

2. Open browser at `http://localhost:8080`

