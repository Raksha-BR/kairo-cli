// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	AI AIConfig `mapstructure:"ai"`
}

type AIConfig struct {
	Provider     string `mapstructure:"provider"`       // "ollama" or "openai"
	OllamaURL    string `mapstructure:"ollama_url"`     // Default: http://localhost:11434
	Model        string `mapstructure:"model"`          // Default: llama3.2
	OpenAIKey    string `mapstructure:"openai_api_key"` // Optional
}

var AppConfig Config

func LoadConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".kairo.yaml")

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set Defaults
	viper.SetDefault("ai.provider", "ollama")
	viper.SetDefault("ai.ollama_url", "http://localhost:11434")
	viper.SetDefault("ai.model", "llama3.2")
	viper.SetDefault("ai.openai_api_key", "")

	// Create default config file if it does not exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := viper.WriteConfigAs(configPath); err != nil {
			return fmt.Errorf("could not create default config file: %v", err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	return viper.Unmarshal(&AppConfig)
}