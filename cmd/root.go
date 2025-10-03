package cmd

import (
	"fmt"
	"os"

	internal "github.com/oreoluwa-bs/zoop/internal"
	"github.com/spf13/cobra"
)

var (
	storeManager *internal.StoreManager
	Version      string
)

var rootCmd = &cobra.Command{
	Use:              "Zoop",
	Short:            "Zoop is a fast cli for storing and retrieving anything",
	Long:             `A fast, minimal CLI for storing and retrieving anything—API keys, passwords, tokens, notes, secrets`,
	PersistentPreRun: preRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, args []string) {
	cfg, err := internal.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var cipher internal.Cipher
	if cfg.Encryption {
		if _, err := os.Stat(cfg.KeyFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Encryption key file not found: %s\n", cfg.KeyFile)
			fmt.Fprintf(os.Stderr, "💡 Run 'zoop init' to generate encryption keys\n")
			os.Exit(1)
		}

		keyData, err := os.ReadFile(cfg.KeyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading key file: %v\n", err)
			os.Exit(1)
		}

		cipher, err = internal.NewAgeCipher(string(keyData))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating cipher: %v\n", err)
			os.Exit(1)
		}
	}

	jsonStore, err := internal.NewJSONStore(cfg.DataFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	store := internal.Store(jsonStore)
	if cipher != nil {
		store = internal.NewEncryptedStore(jsonStore, cipher)
	}

	storeManager = internal.NewStoreManager(store)
}
