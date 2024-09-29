# Project Structure

The folder structure below is based on a three-layer architecture. To gain a clearer understanding of the code and improve readability, it's recommended to first review the [Three-Layer Architecture](Three-Layer-Architecture.md).

Additionally, to understand the process of video streaming on AWS S3, you can refer to the [Video Upload Process](VideoUploadProcess.md). Note that the database only stores the path to the video on AWS S3.


```
.
├── Dockerfile
├── LICENSE.md
├── Makefile
├── README.md
├── assets
├── cmd
│   ├── migration
│   │   └── migration.go
│   └── server
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── i18n
│   ├── de.yaml
│   ├── en.yaml
│   ├── es.yaml
│   ├── fr.yaml
│   ├── it.yaml
│   ├── ja.yaml
│   ├── ko.yaml
│   ├── pt.yaml
│   ├── ru.yaml
│   ├── vi.yaml
│   └── zh.yaml
├── internal
│   ├── entity
│   │   ├── user_entity.go
│   │   └── video_entity.go
│   ├── handler
│   │   └── rest
│   │       └── v1
│   │           ├── handler.go
│   │           ├── user_handler.go
│   │           └── video_handler.go
│   ├── infra
│   │   ├── aws
│   │   │   └── s3.go
│   │   ├── db
│   │   │   ├── database.go
│   │   │   └── redis.go
│   │   ├── env
│   │   │   └── env.go
│   │   ├── reason
│   │   │   └── reason.go
│   │   ├── server
│   │   │   ├── grpc
│   │   │   ├── http
│   │   │   │   ├── http.go
│   │   │   │   └── http_test.go
│   │   │   └── server.go
│   │   └── zap-logging
│   │       ├── log
│   │       │   ├── global.go
│   │       │   ├── level.go
│   │       │   ├── logger.go
│   │       │   └── stdio.go
│   │       └── zap
│   │           ├── option.go
│   │           ├── zap_impl.go
│   │           └── zap_log.go
│   ├── pkg
│   │   ├── json
│   │   │   └── json_handler.go
│   │   ├── localization
│   │   │   └── localization.go
│   │   └── middleware
│   │       ├── auth.go
│   │       └── provider.go
│   ├── repo
│   │   ├── provider.go
│   │   ├── user_repo.go
│   │   └── video_repo.go
│   ├── router
│   │   ├── provider.go
│   │   ├── route.go
│   │   └── swagger_router.go
│   ├── schema
│   │   ├── presigned_url.go
│   │   ├── responses.go
│   │   ├── user_schema.go
│   │   └── video_schema.go
│   └── service
│       ├── auth_service.go
│       ├── provider.go
│       ├── user_service.go
│       └── video_service.go
├── logs
│   ├── mlvt_err_2024-09-07.log
│   └── mlvt_info_2024-09-07.log
├── mlvt.db
├── script
│   ├── build.sh
│   ├── deploy.sh
│   ├── setup.sh
│   └── swagger.sh
└── vendor

34 directories, 73 files
```