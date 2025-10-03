package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
	internal "github.com/oreoluwa-bs/zoop/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize zoop",
	Long:  `Generate encryption keys and set up zoop for first use`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := internal.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		identity, err := age.GenerateX25519Identity()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating key: %v\n", err)
			os.Exit(1)
		}

		if err := os.MkdirAll(filepath.Dir(cfg.KeyFile), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
			os.Exit(1)
		}

		keyFile, err := os.OpenFile(cfg.KeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating key file: %v\n", err)
			os.Exit(1)
		}
		defer keyFile.Close()

		if _, err := keyFile.WriteString(identity.String()); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing key file: %v\n", err)
			os.Exit(1)
		}

		v := viper.New()
		v.SetConfigFile(filepath.Join(filepath.Dir(cfg.KeyFile), "config.yaml"))
		v.Set("encryption", true)
		v.Set("key_file", cfg.KeyFile)
		v.SafeWriteConfig()

		fmt.Printf("✅ Encryption enabled! Key file: %s\n", cfg.KeyFile)
	},
}
