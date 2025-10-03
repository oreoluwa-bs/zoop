package cmd

import (
	"fmt"
	"os"

	internal "github.com/oreoluwa-bs/zoop/internal"
	"github.com/spf13/cobra"
)

var (
	storeManager *internal.StoreManager
)

var rootCmd = &cobra.Command{
	Use:              "Zoop",
	Short:            "Zoop is a fast cli for storing and retrieving you store and retrieve anything",
	Long:             `A fast, minimal CLI for storing and retrieving you store and retrieve anything—API keys, passwords, tokens, notes, secrets`,
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

	jsonStore, err := internal.NewJSONStore(cfg.DataFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var store internal.Store = jsonStore

	if cfg.Encryption {
		if _, err := os.Stat(cfg.KeyFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Encryption key file not found: %s\n", cfg.KeyFile)
			fmt.Fprintf(os.Stderr, "💡 Run 'zoop init' to generate encryption keys\n")
			os.Exit(1)
		}

		cipher, err := internal.NewAgeCipherWithKeyFile(cfg.KeyFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		store = internal.NewEncryptedStore(jsonStore, cipher)
	} else {
		if needsMigration(jsonStore) {
			fmt.Println("⚠️  Detected encrypted data but encryption is disabled")
			fmt.Println("Run 'zoop migrate decrypt' to decrypt your data")
			fmt.Println("Or enable encryption with 'zoop config set encryption true'")
			os.Exit(1)
		}
	}

	storeManager = internal.NewStoreManager(store)
}
