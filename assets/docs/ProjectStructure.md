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
│   └── docs
│       ├── ApiTesting.md
│       ├── AudioFeature.md
│       ├── EnvironmentConfiguration.md
│       ├── FlowEc2.md
│       ├── ProjectStructure.md
│       ├── Three-Layer-Architecture.md
│       ├── TranscriptionFeature.md
│       ├── UserFeature.md
│       ├── VideoFeature.md
│       ├── VideoUploadProcess.md
│       └── swagger.png
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
│   │   ├── audio_entity.go
│   │   ├── frame_entity.go
│   │   ├── transcription_entity.go
│   │   ├── user_entity.go
│   │   └── video_entity.go
│   ├── handler
│   │   └── rest
│   │       └── v1
│   │           ├── audio_handler.go
│   │           ├── handler.go
│   │           ├── transcription_handler.go
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
│   │   ├── audio_repo.go
│   │   ├── provider.go
│   │   ├── transcription_repo.go
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
│       ├── audio_service.go
│       ├── auth_service.go
│       ├── provider.go
│       ├── transcriptions_service.go
│       ├── user_service.go
│       └── video_service.go
├── logs
│   ├── mlvt_err_2024-10-01.log
│   └── mlvt_info_2024-10-01.log
├── mlvt.db
└── script
    ├── build.sh
    ├── deploy.sh
    ├── run_all.sh
    ├── setup.sh
    └── swagger.sh

34 directories, 91 files
```