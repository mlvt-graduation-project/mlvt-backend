package utility

import (
	"fmt"

	"github.com/spf13/viper"
)


func GetPostgresConnection () string {

	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	sslMode := viper.GetString("SSL_MODE")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslMode)
}