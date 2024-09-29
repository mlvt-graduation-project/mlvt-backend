
# Project MLVT

## Table of Contents

- [Project MLVT](#project-mlvt)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Run Script](#run-script)
  - [Getting Started](#getting-started)
    - [Install Dependencies](#install-dependencies)
    - [Run the Server](#run-the-server)
  - [API Documentation](#api-documentation)
  - [Configuration Details](#configuration-details)
  - [Project Architecture](#project-architecture)
  - [API Testing](#api-testing)
  - [Development](#development)
    - [Adding New APIs](#adding-new-apis)
    - [Wire Generation](#wire-generation)
  - [Localization Support](#localization-support)
    - [Supported Languages](#supported-languages)
    - [Changing the Language](#changing-the-language)
  - [Contributing](#contributing)
  - [License](#license)

## Overview

Project MLVT is a robust server application designed to offer a comprehensive suite of APIs for managing users, videos, and additional functionalities. A key feature of this project is its support for generating presigned URLs for storing videos on AWS S3 and managing the video processing pipeline.

This document delivers detailed guidance on setup, architecture, development practices, configuration, execution instructions, API documentation, and localization support.


## Run Script
For those who prefer not to delve deeper into this project, you can simply execute the `make run-all` command. However, please ensure that you install the release version, as other commits may not be stable enough to support the full functionality of the project.


## Getting Started

### Install Dependencies

```bash
make import
# or
go mod tidy
go mod vendor
```

### Run the Server

You can start the server using:

```bash
make run
# or
cd cmd/server
go run .
```

## API Documentation

The API details are available after running the server at `http://localhost:8080/swagger/index.html`. See `docs/swagger.json` for more details.

![Swagger UI](assets/docs/swagger.png)

For API testing instructions, refer to the [API Testing](#api-testing) section.

## Configuration Details

Configuration of the project is managed through environment variables in the `.env` file.

For more details on environment variables, refer to the respective configuration sections under [Environment Configuration](assets/docs/EnvironmentConfiguration.md).

## Project Architecture

* [Three-Layer Architecture](assets/docs/Three-Layer-Architecture.md)
* [Video upload process](assets/docs/VideoUploadProcess.md)
* [Complete project structure](assets/docs/ProjectStructure.md)

## API Testing

* [Core Api Testing](assets/docs/ApiTesting.md)

## Development

### Adding New APIs

When adding new APIs, ensure to add the appropriate annotations before the function. After that, generate the Swagger documentation by running:

```bash
make swag
# or
swag init -g cmd/server/main.go
```

### Wire Generation

To generate the wire files needed for dependency injection, use the following commands:

```bash
make wire
# or
cd cmd/server
wire
```

## Localization Support

This project supports multiple languages for error messages, success notifications, and general information. The localization is implemented using YAML files, stored in the `i18n` directory. Each language has its own YAML file, making it easy to add new languages or update existing translations.

### Supported Languages

- **English** (`en.yaml`): The default language for all messages.
- **Vietnamese** (`vi.yaml`): Translations for Vietnamese users.
- **German** (`de.yaml`): Translations for German users.
- **French** (`fr.yaml`): Translations for French users.
- **Spanish** (`es.yaml`): Translations for Spanish users.
- **Italian** (`it.yaml`): Translations for Italian users.
- **Chinese** (`zh.yaml`): Translations for Chinese users.
- **Japanese** (`ja.yaml`): Translations for Japanese users.
- **Korean** (`ko.yaml`): Translations for Korean users.
- **Portuguese** (`pt.yaml`): Translations for Portuguese users.
- **Russian** (`ru.yaml`): Translations for Russian users.

### Changing the Language

You can change the language of the application by setting the `LANGUAGE` variable in the [Environment Configuration](assets/docs/EnvironmentConfiguration.md). Replace with the appropriate language code from the supported languages list.

## Contributing

We welcome contributions to add more languages, APIs, or improve existing functionalities. Please follow the existing project structure and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for more information.