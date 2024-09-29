
# Environment Configuration

This section outlines the environment variables used in the Project MLVT. 

## How to Setup

1. Copy the template below into a file named `.env` in your project's root directory.
2. Replace the placeholder values with actual configurations suitable for your development or production environment.

## Environment Variables

### Application Settings
```plaintext
APP_NAME=mlvt                      # The name of the application
APP_ENV=development                # Application environment (e.g., development, production)
APP_DEBUG=true                     # Enable debugging (true or false)
```

### Server Configuration
```plaintext
SERVER_PORT=8080                   # The port on which the server will run
```

### Database Configuration
```plaintext
DB_DRIVER=postgres                 # Database driver (e.g., postgres, mysql, sqlite3)
DB_CONNECTION=postgres://username:password@localhost:5432/dbname  # Database connection string
```

### Security Settings
```plaintext
JWT_SECRET=your_secret_key_here    # Secret key for JWT authentication
```

### Logging Configuration
```plaintext
LOG_LEVEL=INFO                    # Set the logging level (INFO, DEBUG, ERROR)
LOG_PATH=./logs/                   # Path where logs are stored
```

### Swagger Configuration
```plaintext
SWAGGER_ENABLED=true               # Enable or disable Swagger documentation (true or false)
SWAGGER_URL=http://localhost:8080/swagger  # URL to access Swagger documentation
```

### AWS S3 Configuration
```plaintext
AWS_REGION=us-west-2               # AWS region for the S3 bucket
AWS_S3_BUCKET=your_bucket_name     # Name of the S3 bucket
AWS_ACCESS_KEY_ID=your_access_key_id           # AWS access key ID
AWS_SECRET_ACCESS_KEY=your_secret_access_key   # AWS secret access key
```

### Language and Localization Settings
```plaintext
LANGUAGE=en                        # Set the language for localization (e.g., en, vi, de)
I18N_PATH=./i18n/                  # Path to the directory containing localization files
```

You can change the language of the application by setting the `LANGUAGE` variable

```env
LANGUAGE="vi"  # For Vietnamese
LANGUAGE="de"  # For German
LANGUAGE="fr"  # For French
LANGUAGE="es"  # For Spanish
LANGUAGE="it"  # For Italian
LANGUAGE="zh"  # For Chinese (Simplified)
LANGUAGE="ja"  # For Japanese
LANGUAGE="ko"  # For Korean
LANGUAGE="pt"  # For Portuguese
LANGUAGE="ru"  # For Russian
```

## Note

- Ensure you do not commit the `.env` file to version control to keep sensitive information like passwords and API keys secure.
- Variables can be adjusted based on specific requirements of different environments (development, staging, production).