# Project Structure

The folder structure below is based on a three-layer architecture. To gain a clearer understanding of the code and improve readability, it's recommended to first review the [Three-Layer Architecture](Three-Layer-Architecture.md).

Additionally, to understand the process of video streaming on AWS S3, you can refer to the [Video Upload Process](VideoUploadProcess.md). Note that the database only stores the path to the video on AWS S3.

```
.
├── assets
│   ├── avatars
│   ├── docs
│   └── videos
├── cmd
│   ├── cleanup
│   ├── migration
│   ├── seeder
│   └── server
├── docs
├── i18n
├── internal
│   ├── entity
│   ├── handler
│   │   └── rest
│   │       └── v1
│   │           ├── audio_handler
│   │           ├── mlvt_handler
│   │           ├── payment_handler
│   │           │   └── momo_handler
│   │           ├── ping_handler
│   │           ├── transcription_handler
│   │           ├── user_handler
│   │           └── video_handler
│   ├── infra
│   │   ├── aws
│   │   ├── db
│   │   │   └── mongodb
│   │   ├── env
│   │   ├── reason
│   │   ├── seeder
│   │   ├── server
│   │   │   ├── grpc
│   │   │   └── http
│   │   └── zap-logging
│   │       ├── log
│   │       └── zap
│   ├── initialize
│   ├── pkg
│   │   ├── json
│   │   ├── localization
│   │   ├── middleware
│   │   ├── request
│   │   └── response
│   ├── repo
│   │   ├── audio_repo
│   │   ├── payment_repo
│   │   │   └── momo_repo
│   │   ├── ping_repo
│   │   ├── transcription_repo
│   │   ├── user_repo
│   │   └── video_repo
│   ├── router
│   ├── schema
│   └── service
│       ├── audio_service
│       ├── auth_service
│       ├── payment_service
│       │   └── momo_service
│       ├── ping_service
│       ├── transcription_service
│       ├── user_service
│       └── video_service
├── logs
├── migration
└── script
```
# Detail:

```
.
├── Dockerfile
├── LICENSE.md
├── Makefile
├── README.md
├── assets
│   ├── avatars
│   │   ├── 11.jpg
│   │   ├── 22.jpg
│   │   ├── 33.jpg
│   │   └── 44.jpg
│   ├── docs
│   │   ├── ApiTesting.md
│   │   ├── AudioFeature.md
│   │   ├── EnvironmentConfiguration.md
│   │   ├── FlowEc2.md
│   │   ├── PingFeature.md
│   │   ├── ProjectStructure.md
│   │   ├── Three-Layer-Architecture.md
│   │   ├── TranscriptionFeature.md
│   │   ├── UserFeature.md
│   │   ├── VideoFeature.md
│   │   ├── VideoUploadProcess.md
│   │   └── swagger.png
│   └── videos
│       ├── 1.mp4
│       ├── 1_thumbnail.jpg
│       ├── 2.mp4
│       ├── 2_thumbnail.jpg
│       ├── 3.mp4
│       ├── 3_thumbnail.jpg
│       ├── 4.mp4
│       └── 4_thumbnail.jpg
├── cmd
│   ├── cleanup
│   │   └── main.go
│   ├── migration
│   │   └── migration.go
│   ├── seeder
│   │   ├── main.go
│   │   └── readme.md
│   └── server
│       └── main.go
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── i18n
│   ├── de.yaml
│   ├── en.yaml
│   ├── es.yaml
│   ├── fr.yaml
│   ├── it.yaml
│   ├── ja.yaml
│   ├── ko.yaml
│   ├── pt.yaml
│   ├── ru.yaml
│   ├── vi.yaml
│   └── zh.yaml
├── internal
│   ├── entity
│   │   ├── audio_entity.go
│   │   ├── constants.go
│   │   ├── frame_entity.go
│   │   ├── momo_payment_entity.go
│   │   ├── transaction_entity.go
│   │   ├── transcription_entity.go
│   │   ├── user_entity.go
│   │   └── video_entity.go
│   ├── handler
│   │   └── rest
│   │       └── v1
│   │           ├── audio_handler
│   │           │   └── audio_handler.go
│   │           ├── handler.go
│   │           ├── mlvt_handler
│   │           │   └── mlvt_handler.go
│   │           ├── payment_handler
│   │           │   └── momo_handler
│   │           │       └── momo_payment_handler.go
│   │           ├── ping_handler
│   │           │   └── ping_handler.go
│   │           ├── transcription_handler
│   │           │   └── transcription_handler.go
│   │           ├── user_handler
│   │           │   ├── user_handler.go
│   │           │   ├── user_handler_mock.go
│   │           │   └── user_handler_test.go
│   │           └── video_handler
│   │               ├── video_handler.go
│   │               └── video_handler_test.go
│   ├── infra
│   │   ├── aws
│   │   │   ├── provider.go
│   │   │   ├── s3.go
│   │   │   └── s3_mock.go
│   │   ├── db
│   │   │   ├── database.go
│   │   │   ├── mongodb
│   │   │   │   ├── filter.go
│   │   │   │   ├── mongodb_adapter.go
│   │   │   │   └── mongodb_client.go
│   │   │   └── redis.go
│   │   ├── env
│   │   │   └── env.go
│   │   ├── reason
│   │   │   └── reason.go
│   │   ├── seeder
│   │   │   ├── user_seeder.go
│   │   │   └── user_video_seeder.go
│   │   ├── server
│   │   │   ├── grpc
│   │   │   ├── http
│   │   │   │   ├── http.go
│   │   │   │   └── http_test.go
│   │   │   └── server.go
│   │   └── zap-logging
│   │       ├── log
│   │       │   ├── global.go
│   │       │   ├── level.go
│   │       │   ├── logger.go
│   │       │   └── stdio.go
│   │       └── zap
│   │           ├── option.go
│   │           ├── zap_impl.go
│   │           └── zap_log.go
│   ├── initialize
│   │   ├── aws.go
│   │   ├── database.go
│   │   ├── logger.go
│   │   ├── router.go
│   │   ├── run.go
│   │   ├── server.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── pkg
│   │   ├── json
│   │   │   └── json_handler.go
│   │   ├── localization
│   │   │   └── localization.go
│   │   ├── middleware
│   │   │   ├── auth.go
│   │   │   ├── auth_mock.go
│   │   │   └── provider.go
│   │   ├── request
│   │   │   ├── base_request.go
│   │   │   ├── ls_request.go
│   │   │   ├── stt_request.go
│   │   │   ├── tts_request.go
│   │   │   └── ttt_request.go
│   │   └── response
│   │       ├── ec2_response.go
│   │       └── response.go
│   ├── repo
│   │   ├── audio_repo
│   │   │   └── audio_repo.go
│   │   ├── payment_repo
│   │   │   └── momo_repo
│   │   │       └── momo_payment_repo.go
│   │   ├── ping_repo
│   │   │   └── ping_repo.go
│   │   ├── provider.go
│   │   ├── transaction_log_repo.go
│   │   ├── transcription_repo
│   │   │   └── transcription_repo.go
│   │   ├── user_repo
│   │   │   ├── user_repo.go
│   │   │   ├── user_repo_mock.go
│   │   │   └── user_repo_test.go
│   │   └── video_repo
│   │       ├── video_repo.go
│   │       ├── video_repo_mock.go
│   │       └── video_repo_test.go
│   ├── router
│   │   ├── provider.go
│   │   ├── route.go
│   │   └── swagger_router.go
│   ├── schema
│   │   ├── presigned_url.go
│   │   ├── responses.go
│   │   ├── user_schema.go
│   │   └── video_schema.go
│   └── service
│       ├── audio_service
│       │   └── audio_service.go
│       ├── auth_service
│       │   ├── auth_service.go
│       │   └── auth_service_mock.go
│       ├── payment_service
│       │   └── momo_service
│       │       └── momo_payment_service.go
│       ├── ping_service
│       │   └── ping_service.go
│       ├── provider.go
│       ├── transcription_service
│       │   └── transcriptions_service.go
│       ├── user_service
│       │   ├── user_service.go
│       │   ├── user_service_mock.go
│       │   └── user_service_test.go
│       └── video_service
│           ├── video_service.go
│           ├── video_service_mock.go
│           └── video_service_test.go
├── logs
│   ├── mlvt_err_2024-11-28.log
│   ├── mlvt_err_2024-12-01.log
│   ├── mlvt_err_2024-12-04.log
│   ├── mlvt_err_2024-12-05.log
│   ├── mlvt_info_2024-11-28.log
│   ├── mlvt_info_2024-12-01.log
│   ├── mlvt_info_2024-12-04.log
│   └── mlvt_info_2024-12-05.log
├── migration
│   ├── 0001_create_users_table.down.sql
│   ├── 0001_create_users_table.up.sql
│   ├── 0002_create_videos_table.down.sql
│   ├── 0002_create_videos_table.up.sql
│   ├── 0003_create_transcriptions_table.down.sql
│   ├── 0003_create_transcriptions_table.up.sql
│   ├── 0004_create_transaction_logs_table.down.sql
│   ├── 0004_create_transaction_logs_table.up.sql
│   ├── 0005_create_frames_table.down.sql
│   ├── 0005_create_frames_table.up.sql
│   ├── 0006_create_audios_table.down.sql
│   ├── 0006_create_audios_table.up.sql
│   ├── 0007_insert_sample_users.down.sql
│   └── 0007_insert_sample_users.up.sql
├── mlvt.db
└── script
    ├── build.sh
    ├── deploy.sh
    ├── run_all.sh
    ├── setup.sh
    └── swagger.sh
```