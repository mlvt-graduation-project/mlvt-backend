package env

import (
	"fmt"
	"path/filepath"
	"sync"

	"mlvt/internal/infra/zap-logging/log"

	"github.com/spf13/viper"
)

var (
	EnvConfig *Config // Singleton instance of the configuration
	mu        sync.RWMutex
)

const defaultEnvFilePath = ".env"

// Config holds all the environment variables used in the application.
type Config struct {
	AppName        string
	AppEnv         string
	AppDebug       bool
	ServerPort     string
	LogLevel       string
	LogPath        string
	DBDriver       string
	DBConnection   string
	JWTSecret      string
	SwaggerEnabled bool
	SwaggerURL     string
	AWSRegion      string
	AWSBucket      string
	AWSAccessKeyID string
	AWSSecretKey   string
	Language       string
	I18NPath       string
	RootDir        string
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
	dbPath := resolvePath(rootDir, viper.GetString("DB_CONNECTION"))

	EnvConfig = &Config{
		AppName:        viper.GetString("APP_NAME"),
		AppEnv:         viper.GetString("APP_ENV"),
		AppDebug:       viper.GetBool("APP_DEBUG"),
		ServerPort:     viper.GetString("SERVER_PORT"),
		LogLevel:       viper.GetString("LOG_LEVEL"),
		LogPath:        logPath,
		DBDriver:       viper.GetString("DB_DRIVER"),
		DBConnection:   dbPath,
		JWTSecret:      viper.GetString("JWT_SECRET"),
		SwaggerEnabled: viper.GetBool("SWAGGER_ENABLED"),
		SwaggerURL:     viper.GetString("SWAGGER_URL"),
		AWSRegion:      viper.GetString("AWS_REGION"),
		AWSBucket:      viper.GetString("AWS_BUCKET"),
		AWSAccessKeyID: viper.GetString("AWS_ACCESS_KEY_ID"),
		AWSSecretKey:   viper.GetString("AWS_SECRET_KEY"),
		Language:       viper.GetString("LANGUAGE"),
		I18NPath:       i18nPath,
		RootDir:        rootDir,
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
	rootDir := "/Users/giabao/Code/golang/mlvt/mlvt-backend"
	return filepath.Abs(rootDir)
}

// resolvePath combines the root directory with a relative path.
func resolvePath(rootDir, relPath string) string {
	if filepath.IsAbs(relPath) {
		return relPath
	}
	return filepath.Join(rootDir, relPath)
}
