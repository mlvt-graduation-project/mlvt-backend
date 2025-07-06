
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

### VietQR Payment Configuration
```plaintext
# VietQR API credentials (get from https://my.vietqr.io)
VIETQR_CLIENT_ID=your_client_id_here      # VietQR Client ID
VIETQR_API_KEY=your_api_key_here          # VietQR API Key

# VietinBank account information for receiving payments
VIETINBANK_ACCOUNT_NO=your_account_number    # Your VietinBank account number (6-19 digits)
VIETINBANK_ACCOUNT_NAME=YOUR ACCOUNT NAME    # Account holder name (uppercase, no special chars)
VIETINBANK_BIN_CODE=970415                   # VietinBank BIN code (always 970415)
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

## Payment Options

The QR payment feature supports the following payment packages:

| Option | Tokens | VND Price | Description |
|--------|--------|-----------|-------------|
| 5k   | 500    | 5,000   | 500 tokens - 5,000 VND |
| 10k     | 1,000  | 10,000 | 1,000 tokens - 10,000 VND |
| 20k     | 2,000  | 20,000 | 2,000 tokens - 20,000 VND |
| 50k     | 5,000  | 50,000 | 5,000 tokens - 50,000 VND |
| 100k    | 10,000 | 100,000| 10,000 tokens - 100,000 VND |

## Note

- Ensure you do not commit the `.env` file to version control to keep sensitive information like passwords and API keys secure.
- Variables can be adjusted based on specific requirements of different environments (development, staging, production).
- For VietQR API access, register at [My VietQR](https://my.vietqr.io) to obtain your Client ID and API Key.
- The VietinBank BIN code is always `970415` for VietinBank.
