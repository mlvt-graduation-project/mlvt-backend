package reason

import "mlvt/internal/pkg/localization"

// Define the reason constants using the LocalizedString type
var (
	// Error messages under 'error.user'
	InvalidRequest           localization.LocalizedString = "error.user.invalid_request"
	InvalidRequestFormat     localization.LocalizedString = "error.user.invalid_request_format"
	InternalServerError      localization.LocalizedString = "error.user.internal_server_error"
	InvalidUserID            localization.LocalizedString = "error.user.invalid_user_id"
	UserNotFound             localization.LocalizedString = "error.user.not_found"
	Unauthorized             localization.LocalizedString = "error.user.unauthorized"
	EmailAlreadyRegistered   localization.LocalizedString = "error.user.email_already_registered"
	InvalidCredentials       localization.LocalizedString = "error.user.invalid_credentials"
	FailedToGenerateToken    localization.LocalizedString = "error.user.failed_to_generate_token"
	UnexpectedSigningMethod  localization.LocalizedString = "error.user.unexpected_signing_method"
	InvalidToken             localization.LocalizedString = "error.user.invalid_token"
	InvalidTokenClaims       localization.LocalizedString = "error.user.invalid_token_claims"
	InvalidUserIDTypeInToken localization.LocalizedString = "error.user.invalid_userid_type_in_token"

	// Error messages under 'error.video'
	VideoInvalidRequest         localization.LocalizedString = "error.video.invalid_request"
	VideoInternalServerError    localization.LocalizedString = "error.video.internal_server_error"
	InvalidVideoID              localization.LocalizedString = "error.video.invalid_video_id"
	VideoNotFound               localization.LocalizedString = "error.video.not_found"
	FailedToAddVideo            localization.LocalizedString = "error.video.failed_to_add_video"
	VideoLinkCannotBeEmpty      localization.LocalizedString = "error.video.video_link_cannot_be_empty"
	VideoDurationMustBePositive localization.LocalizedString = "error.video.video_duration_must_be_positive"
	VideoTitleCannotBeEmpty     localization.LocalizedString = "error.video.video_title_cannot_be_empty"
	NoVideoForUser              localization.LocalizedString = "error.video.no_video_for_user"

	// Error messages under 'error.data'
	InsertSampleFailed              localization.LocalizedString = "error.data.insert_sample"
	MigrationFailed                 localization.LocalizedString = "error.data.migration_failed"
	FailedToInitializeAWSS3         localization.LocalizedString = "error.data.failed_to_initialize_aws_s3"
	FailedToInitializeDB            localization.LocalizedString = "error.data.failed_to_initialize_db"
	FailedToCreatePresignedURL      localization.LocalizedString = "error.data.failed_to_create_presigned_url"
	FailedToGeneratePresignedURL    localization.LocalizedString = "error.data.failed_to_generate_presigned_url"
	FailedToPresignPutObjectRequest localization.LocalizedString = "error.data.failed_to_presign_put_object_request"
	RequestFormatError              localization.LocalizedString = "error.data.request_format_error"
	ErrorReadingYAMLFile            localization.LocalizedString = "error.data.error_reading_yaml"
	ErrorUnmarshalingYAMLData       localization.LocalizedString = "error.data.error_unmarshaling_yaml"

	// Error messages under 'error.env'
	UnableToLoadAWSConfig localization.LocalizedString = "error.env.unable_to_load_aws_config"
	ErrorLoadingEnvFile   localization.LocalizedString = "error.env.error_loading_env"
	I18nPathNotSet        localization.LocalizedString = "error.env.i18n_path_not_set"

	// Error messages under 'error.general'
	ServerStartFailed         localization.LocalizedString = "error.general.server_start_failed"
	ErrorOccurred             localization.LocalizedString = "error.general.error_occurred"
	FailedToBindJSON          localization.LocalizedString = "error.general.failed_to_bind_json"
	KeyNotFoundOrTypeMismatch localization.LocalizedString = "error.general.key_not_found_or_type_mismatch"
	KeyNotFound               localization.LocalizedString = "error.general.key_not_found"
	MessageNotFound           localization.LocalizedString = "error.general.message_not_found"

	// Success messages under 'success.user'
	UserRegistered localization.LocalizedString = "success.user.registered"
	UserUpdated    localization.LocalizedString = "success.user.updated"
	UserDeleted    localization.LocalizedString = "success.user.deleted"

	// Success messages under 'success.video'
	VideoAdded   localization.LocalizedString = "success.video.added"
	VideoUpdated localization.LocalizedString = "success.video.updated"
	VideoDeleted localization.LocalizedString = "success.video.deleted"

	// Success messages under 'success.general'
	MigrationsApplied             localization.LocalizedString = "success.general.migrations_applied"
	ServerStarted                 localization.LocalizedString = "success.general.server_start"
	JSONRequestBodyRead           localization.LocalizedString = "success.general.json_request_body_read"
	DatabaseConnectionEstablished localization.LocalizedString = "success.general.database_connection_established"
	InsertSampleDataSuccess       localization.LocalizedString = "success.general.insert_sample"

	// Common messages under 'common.info'
	Error                        localization.LocalizedString = "common.info.error"
	Status                       localization.LocalizedString = "common.info.status"
	Data                         localization.LocalizedString = "common.info.data"
	Type                         localization.LocalizedString = "common.info.type"
	GeneratedPresignedURL        localization.LocalizedString = "common.info.generated_presigned_url"
	GeneratedPresignedURLForFile localization.LocalizedString = "common.info.generated_presigned_url_for_file"
	Migration                    localization.LocalizedString = "common.info.migration"
	AlreadyApplied               localization.LocalizedString = "common.info.already_applied"
	ApplyingMigration            localization.LocalizedString = "common.info.applying_migration"
	ResponseWritten              localization.LocalizedString = "common.info.response_written"
	LoadedMessagesForLanguage    localization.LocalizedString = "common.info.loaded_messages_for_language"
	ServerShutdown               localization.LocalizedString = "common.info.server_shutdown"
	ServerForcedShutdown         localization.LocalizedString = "common.info.server_forced_shutdown"
)
