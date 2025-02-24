package entity

type TrafficActionType string

const (
	LoginAction         TrafficActionType = "login"
	LogoutAction        TrafficActionType = "logout"
	CreateAccountAction TrafficActionType = "create_account"

	ProcessSTTModelAction          TrafficActionType = "process_stt_model"
	ProcessTTTModelAction          TrafficActionType = "process_ttt_model"
	ProcessTTSModelAction          TrafficActionType = "process_tts_model"
	ProcessLSModelAction           TrafficActionType = "process_ls_model"
	ProcessFullPipelineModelAction TrafficActionType = "process_full_pipeline"

	// related to progress
	ProgressDurationAction TrafficActionType = "progress_duration"
	ProgressStatusAction   TrafficActionType = "progress_status"

	// related to profile
	UploadAvatarAction   TrafficActionType = "upload_avatar"
	UpdateProfileAction  TrafficActionType = "update_profile"
	ChangePasswordAction TrafficActionType = "change_password"

	// related to user wallet
	UseTokenAction      TrafficActionType = "use_token"
	RedeemVoucherAction TrafficActionType = "redeem_voucher"

	// file handling
	UploadFileAction TrafficActionType = "upload_file"
	DeleteFileAction TrafficActionType = "delete_file"

	// activities belong to admin
	AdminUserProfileAction      TrafficActionType = "admin_user_profile" // edit: info, active-inactive
	AdminVoucherAction          TrafficActionType = "admin_voucher"      // CRUD
	AdminModelOptionAction      TrafficActionType = "admin_model_option" // CRUD
	AdminConfigAction           TrafficActionType = "admin_config"       // CRUD
	AdminSetDefaultConfigAction TrafficActionType = "admin_set_default_config"
)
