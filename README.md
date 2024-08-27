
# Project MLVT

## Overview

This project is structured to provide a robust server application with various APIs. The documentation and running instructions are outlined below.

## API Documentation

You can find the API documentation in the following file:
- `docs/swagger.json`

## Getting Started

To set up and run the project, follow these steps:

1. **Install dependencies**  
   Run the following command to tidy up and install Go module dependencies:
   ```bash
   go mod tidy
   ```

2. **Run the server**  
   You can run the server using either of the following methods:
   ```bash
   make run
   ```
   or
   ```bash
   cd cmd/server
   go run .
   ```

## Wire Generation

To generate the wire files needed for dependency injection, you can use the following commands:

```bash
make wire
```
or
```bash
cd cmd/server
wire
```

## Adding New APIs

When adding new APIs, ensure to add the appropriate annotations before the function. After that, generate the Swagger documentation by running:

```bash
swag init -g cmd/server/main.go
```

## Project Structure

Below is the current folder structure of the project:

```
.
├── cmd
│   ├── migration
│   └── server
├── docs
├── i18n
├── internal
│   ├── entity
│   ├── handler
│   │   └── rest
│   │       └── v1
│   ├── infra
│   │   ├── cli
│   │   ├── conf
│   │   │   └── viper
│   │   │       └── testdata
│   │   ├── db
│   │   ├── server
│   │   │   ├── grpc
│   │   │   └── http
│   │   ├── translator
│   │   ├── validator
│   │   └── zap-logging
│   │       ├── log
│   │       └── zap
│   ├── pkg
│   │   ├── json
│   │   └── middleware
│   ├── repo
│   ├── router
│   ├── schema
│   ├── service
│   └── translator
├── logs
└── script

35 directories
```

