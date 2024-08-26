package viper

import (
	"fmt"
	"mlvt/internal/infra/conf"
	"os"

	"github.com/spf13/viper"
)

type FileConfig struct {
	Path   string `json:"path"`
	parser *viper.Viper
}

// NewWithPath to create config from specified path
func NewWithPath(filePath string) (c conf.Config, err error) {
	var stat os.FileInfo

	stat, err = os.Stat(filePath)
	if err != nil {
		return
	}

	if !stat.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", filePath)
	}

	p := viper.New()
	p.SetConfigFile(filePath)

	err = p.ReadInConfig()
	if err != nil {
		return
	}

	return &FileConfig{Path: filePath, parser: p}, nil
}

// Parse parses the configuration by object pointer
func (c *FileConfig) Parse(obj any) error {
	return c.parser.Unmarshal(obj)
}

func (c *FileConfig) Get(key string) any {
	return c.parser.Get(key)
}

func (c *FileConfig) GetBool(key string) bool {
	return c.parser.GetBool(key)
}

func (c *FileConfig) GetString(s string) string {
	return c.parser.GetString(s)
}

func (c *FileConfig) GetInt(s string) int {
	return c.parser.GetInt(s)
}
