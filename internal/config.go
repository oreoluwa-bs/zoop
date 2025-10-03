package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DataFile string `mapstructure:"data_file"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// Expand $HOME before setting default
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	v.SetDefault("data_file", filepath.Join(homeDir, ".zoop", "store.json"))

	// v.SetDefault("encryption", false)

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME/.zoop")
	v.AddConfigPath(".")

	v.AutomaticEnv()
	v.SetEnvPrefix("ZOOP")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found is OK, use defaults
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(config.DataFile), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	// v.SafeWriteConfig()

	return &config, nil
}
