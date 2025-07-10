package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mlvt/internal/entity"
	"mlvt/internal/schema"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func GetPostgresConnection() string {
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	sslMode := viper.GetString("SSL_MODE")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslMode)
}

func SortProcessResponseByCreateDate(p []entity.Process, limit int) []entity.Process {
	// Sort by created at in descending order
	sort.Slice(p, func(i, j int) bool {
		return p[i].CreatedAt.After(p[j].CreatedAt)
	})

	// Limit the result
	if limit > 0 && limit < len(p) {
		return p[:limit]
	}
	return p
}

func Contains[T comparable](list []T, target T) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func GetMediaTitle(mediaType entity.MediaType, isFullPipeline bool, isOriginalText bool, id int) string {
	var result string
	switch mediaType {
	case entity.MediaTypeAudio:
		result = fmt.Sprintf("Audio - %v", id)
	case entity.MediaTypeVideo:
		result = fmt.Sprintf("Video - %v", id)
	case entity.MediaTypeText:
		if isOriginalText {
			result = fmt.Sprintf("Text - %v", id)
		} else {
			result = fmt.Sprintf("Translated Text - %v", id)
		}
	}
	if isFullPipeline {
		result += " (Full pipeline)"
	}
	return result
}

var secretKey = []byte("my32byteSuperSecretKey1234567890")

func SetExpireTime(minutes int) time.Time {
	now := time.Now()
	expire := now.Add(time.Duration(minutes) * time.Minute)
	return expire
}

// EncryptToken encrypts username + expire_date to a base64 string
func EncryptToken(username string, expireDate time.Time) (string, error) {
	payload := schema.TokenPayload{
		Username:   username,
		ExpireDate: expireDate,
	}

	// Chuyển về JSON
	plainData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	// GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, plainData, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptToken reverses the encrypted token to get username + expire_date
func DecryptToken(token string) (string, time.Time, error) {
	cipherData, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", time.Time{}, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", time.Time{}, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherData) < nonceSize {
		return "", time.Time{}, fmt.Errorf("invalid token")
	}

	nonce, ciphertext := cipherData[:nonceSize], cipherData[nonceSize:]

	plainData, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	var payload schema.TokenPayload
	if err := json.Unmarshal(plainData, &payload); err != nil {
		return "", time.Time{}, err
	}

	return payload.Username, payload.ExpireDate, nil
}

func IsInListString(value string, list []string) bool {
	for _, item := range list {
		if strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}

func GetProgressTitle(id int, progressType entity.ProgressType) string {
	var prefixTitle string

	switch progressType {
	case entity.ProgressTypeSTT:
		prefixTitle = "Text Generation"
	case entity.ProgressTypeTTT:
		prefixTitle = "Text Translation"
	case entity.ProgressTypeTTS:
		prefixTitle = "Voice Generation"
	case entity.ProgressTypeLS:
		prefixTitle = "Lip Synchronization"
	case entity.ProgressTypeFP:
		prefixTitle = "Video Translation"
	}

	return fmt.Sprintf("%v - %v", prefixTitle, id)
}
