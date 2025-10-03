package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage zoop configuration",
	Long:  `Set, get, or list configuration values persistently`,
}

var configSetCmd = &cobra.Command{
	Use:   "set key=value",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Error: Invalid format. Use key=value\n")
			os.Exit(1)
		}

		key, value := parts[0], parts[1]

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		configDir := filepath.Join(homeDir, ".zoop")
		configPath := filepath.Join(configDir, "config.yaml")

		// Load existing config as map
		config := make(map[string]interface{})
		if data, err := os.ReadFile(configPath); err == nil {
			if err := yaml.Unmarshal(data, &config); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing existing config: %v\n", err)
				os.Exit(1)
			}
		}

		// Set the value (parse bool, etc.)
		if value == "true" {
			config[key] = true
		} else if value == "false" {
			config[key] = false
		} else {
			config[key] = value
		}

		// Ensure directory exists
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			os.Exit(1)
		}

		// Write back
		data, err := yaml.Marshal(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling config: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Set %s = %s\n", key, value)
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get key",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		configDir := filepath.Join(homeDir, ".zoop")

		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(configDir)

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Fprintf(os.Stderr, "Error: Config file not found\n")
			} else {
				fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			}
			os.Exit(1)
		}

		value := v.Get(key)
		if value == nil {
			fmt.Fprintf(os.Stderr, "Error: Config key '%s' not found\n", key)
			os.Exit(1)
		}
		fmt.Println(value)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		configDir := filepath.Join(homeDir, ".zoop")

		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(configDir)

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("No config file found")
			} else {
				fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			}
			return
		}

		settings := v.AllSettings()
		for key, value := range settings {
			fmt.Printf("%s: %v\n", key, value)
		}
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
}
