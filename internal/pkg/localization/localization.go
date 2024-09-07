package localization

import (
	"fmt"
	"io/ioutil"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	messages map[string]interface{} // Map to store loaded messages
	mu       sync.RWMutex           // Mutex to handle concurrent access
)

const defaultLanguage = "en"

// LocalizedString represents a localized string key
type LocalizedString string

// init loads the default language file (e.g., English) at startup
func init() {
	// Get the language and path from environment variables
	lang := env.EnvConfig.Language
	if lang == "" {
		lang = defaultLanguage // Use default language if not set
	}

	i18nPath := env.EnvConfig.I18NPath
	if i18nPath == "" {
		log.Error("I18N_PATH environment variable is not set")
	}

	if i18nPath != "" && i18nPath[len(i18nPath)-1] != '/' {
		i18nPath += "/"
	}

	// Load the YAML file for the specified language
	loadYAML(lang, i18nPath)
}

// loadYAML reads and unmarshals the YAML file for the specified language
func loadYAML(lang string, path string) {
	mu.Lock()
	defer mu.Unlock()

	filePath := fmt.Sprintf("%s%s.yaml", path, lang)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Error reading YAML file: %v", err)
	}

	var yamlData map[string]interface{}
	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		log.Errorf("Error unmarshaling YAML data: %v", err)
	}

	messages = yamlData

	// Debugging: Print out the loaded messages map
	fmt.Printf("Loaded messages for language '%s'\n", lang)
}

// SetLanguage switches the localization to the specified language
func SetLanguage(lang string) {
	i18nPath := env.EnvConfig.I18NPath
	if i18nPath == "" {
		log.Error("I18N_PATH environment variable is not set")
	}
	loadYAML(lang, i18nPath)
}

// Message retrieves the localized message using the key defined in LocalizedString
func (ls LocalizedString) Message() string {
	mu.RLock()
	defer mu.RUnlock()

	keys := strings.Split(string(ls), ".")
	var result interface{} = messages

	for _, key := range keys {
		switch value := result.(type) {
		case map[interface{}]interface{}:
			result = value[key]
		case map[string]interface{}:
			result = value[key]
		default:
			fmt.Printf("Key not found or type mismatch: %s\n", key)
			return "Message not found"
		}

		if result == nil {
			fmt.Printf("Key not found: %s\n", key)
			return "Message not found"
		}
	}

	if msg, ok := result.(string); ok {
		return msg
	}

	return "Message not found"
}
