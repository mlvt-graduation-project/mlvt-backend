package env

import (
	"fmt"
	"path/filepath"
	"sync"

	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/utility"

	"github.com/spf13/viper"
)

var (
	EnvConfig *Config // Singleton instance of the configuration
	mu        sync.RWMutex
)

const defaultEnvFilePath = ".env"

// Config holds all the environment variables used in the application.
type Config struct {
	AppName    string
	AppEnv     string
	AppDebug   bool
	ServerPort string
	LogLevel   string
	LogPath    string

	// db struct
	DBDriver     string
	DBConnection string

	JWTSecret             string
	SwaggerEnabled        bool
	SwaggerURL            string
	AWSRegion             string
	AWSBucket             string
	AWSAccessKeyID        string
	AWSSecretKey          string
	AudioFolder           string
	AvatarFolder          string
	VideosFolder          string
	TranscriptionsFolder  string
	VideoFramesFolder     string
	Ec2IPAddress          string
	Ec2Port               string
	Language              string
	MigrationsPath        string
	MongoDBEndPoint       string
	I18NPath              string
	RootDir               string
	TelegramBotToken      string
	TelegramChatID        string
	VietQRClientID        string
	VietQRAPIKey          string
	VietinBankAccountNo   string
	VietinBankAccountName string
	VietinBankBinCode     string
	SMTPEmail             string
	SMTPPassword          string
	SMTPHost              string
	SMTPPort              string
}

// init loads the environment variables at startup
func init() {
	mu.Lock()
	defer mu.Unlock()

	// Load environment variables from .env file
	err := loadEnvFile(defaultEnvFilePath)
	if err != nil {
		log.Errorf("Error loading environment variables from file: %v", err)
	}

	// Initialize EnvConfig
	err = initializeConfig()
	if err != nil {
		log.Errorf("Error initializing configuration: %v", err)
	}
}

// loadEnvFile loads environment variables from the specified .env file
func loadEnvFile(filePath string) error {
	rootDir, err := getProjectRootDir()
	if err != nil {
		return fmt.Errorf("error getting project root directory: %v", err)
	}

	viper.SetConfigFile(filepath.Join(rootDir, filePath))
	viper.AutomaticEnv()

	// Load the environment file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}
	return nil
}

// initializeConfig sets up the EnvConfig singleton from environment variables
func initializeConfig() error {
	rootDir, err := getProjectRootDir()
	if err != nil {
		return err
	}

	// Adjust relative paths
	logPath := resolvePath(rootDir, viper.GetString("LOG_PATH"))
	i18nPath := resolvePath(rootDir, viper.GetString("I18N_PATH"))
	dbPath := utility.GetPostgresConnection()

	EnvConfig = &Config{
		AppName:               viper.GetString("APP_NAME"),
		AppEnv:                viper.GetString("APP_ENV"),
		AppDebug:              viper.GetBool("APP_DEBUG"),
		ServerPort:            viper.GetString("SERVER_PORT"),
		LogLevel:              viper.GetString("LOG_LEVEL"),
		LogPath:               logPath,
		DBDriver:              viper.GetString("DB_DRIVER"),
		DBConnection:          dbPath,
		JWTSecret:             viper.GetString("JWT_SECRET"),
		SwaggerEnabled:        viper.GetBool("SWAGGER_ENABLED"),
		SwaggerURL:            viper.GetString("SWAGGER_URL"),
		AWSRegion:             viper.GetString("AWS_REGION"),
		AWSBucket:             viper.GetString("AWS_BUCKET"),
		AWSAccessKeyID:        viper.GetString("AWS_ACCESS_KEY_ID"),
		AWSSecretKey:          viper.GetString("AWS_SECRET_KEY"),
		Language:              viper.GetString("LANGUAGE"),
		AudioFolder:           viper.GetString("AUDIO_FOLDER"),
		AvatarFolder:          viper.GetString("AVATAR_FOLDER"),
		VideosFolder:          viper.GetString("VIDEOS_FOLDER"),
		TranscriptionsFolder:  viper.GetString("TRANSCRIPTIONS_FOLDER"),
		VideoFramesFolder:     viper.GetString("VIDEO_FRAMES_FOLDER"),
		Ec2IPAddress:          viper.GetString("EC2_IP_ADDRESS"),
		Ec2Port:               viper.GetString("EC2_PORT"),
		MigrationsPath:        viper.GetString("MIGRATIONS_PATH"),
		MongoDBEndPoint:       viper.GetString("MONGODB_ENDPOINT"),
		SMTPEmail:             viper.GetString("SENDER_EMAIL"),
		SMTPPassword:          viper.GetString("SMTP_PASSWORD"),
		SMTPHost:              viper.GetString("SMTP_HOST"),
		SMTPPort:              viper.GetString("SMTP_PORT"),
		I18NPath:              i18nPath,
		RootDir:               rootDir,
		TelegramBotToken:      viper.GetString("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:        viper.GetString("TELEGRAM_CHAT_ID"),
		VietQRClientID:        viper.GetString("VIETQR_CLIENT_ID"),
		VietQRAPIKey:          viper.GetString("VIETQR_API_KEY"),
		VietinBankAccountNo:   viper.GetString("VIETINBANK_ACCOUNT_NO"),
		VietinBankAccountName: viper.GetString("VIETINBANK_ACCOUNT_NAME"),
		VietinBankBinCode:     viper.GetString("VIETINBANK_BIN_CODE"),
	}

	if EnvConfig.JWTSecret == "" {
		return fmt.Errorf("required environment variable JWT_SECRET is not set")
	}
	return nil
}

// getProjectRootDir returns the root directory of the project
func getProjectRootDir() (string, error) {
	// Set your own root go manage the .env
	//Ex: rootDir := "/Users/giabao/Code/golang/mlvt/mlvt-backend"
	rootDir := "../../"
	return filepath.Abs(rootDir)
}

// resolvePath combines the root directory with a relative path.
func resolvePath(rootDir, relPath string) string {
	if filepath.IsAbs(relPath) {
		return relPath
	}
	return filepath.Join(rootDir, relPath)
}
