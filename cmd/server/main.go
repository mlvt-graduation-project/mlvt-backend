package main

import (
	"context"
	"flag"
	"log"
	infra "mlvt/internal/infra/database"
	"mlvt/internal/pkg/auth"
	"mlvt/internal/routes"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

type application struct {
	DSN          string
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	Domain       string
	APIKey       string
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var app application

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.StringVar(&app.APIKey, "api-key", "", "api key")
	flag.Parse()

	infra.ConnectDatabase()
	db := infra.DB
	defer db.Close()

	authConfig := &auth.Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Hour * 24,
		RefreshExpiry: time.Hour * 24 * 7,
		CookieDomain:  app.CookieDomain,
		CookiePath:    "/",
		CookieName:    "refresh_token",
	}

	// AWS S3 client setup
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	router := routes.SetupRouter(db, s3Client, authConfig)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
