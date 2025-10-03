package cmd

import (
	"fmt"
	"os"
	"strings"

	internal "github.com/oreoluwa-bs/zoop/internal"
	"github.com/spf13/cobra"
)

// cmd/migrate.go
var migrateCmd = &cobra.Command{
	Use:   "migrate [decrypt|encrypt]",
	Short: "Migrate data between encrypted and unencrypted storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]

		cfg, err := internal.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if action == "decrypt" {
			cipher, err := internal.NewAgeCipherWithKeyFile(cfg.KeyFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			encInternalStore, err := internal.NewJSONStore(cfg.DataFile + ".enc")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			encryptedStore := internal.NewEncryptedStore(
				encInternalStore, cipher)
			plainStore, err := internal.NewJSONStore(cfg.DataFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Migrate all data
			err = internal.MigrateStores(encryptedStore, plainStore)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Clean up encrypted file
			err = os.Remove(cfg.DataFile + ".enc")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to remove encrypted file: %v\n", err)
			}

			fmt.Println("✅ Data decrypted successfully")

		} else if action == "encrypt" {
			cipher, err := internal.NewAgeCipherWithKeyFile(cfg.KeyFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			encInternalStore, err := internal.NewJSONStore(cfg.DataFile + ".enc")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			encryptedStore := internal.NewEncryptedStore(
				encInternalStore, cipher)
			plainStore, err := internal.NewJSONStore(cfg.DataFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Migrate all data
			err = internal.MigrateStores(plainStore, encryptedStore)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Clean up plain file
			err = os.Remove(cfg.DataFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to remove plain file: %v\n", err)
			}

			fmt.Println("✅ Data encrypted successfully")
		}
	},
}

func needsMigration(store internal.Store) bool {
	keys, err := store.GetAllKeys()
	if err != nil || len(keys) == 0 {
		return false
	}
	sampleValue, err := store.Get(keys[0])
	if err != nil {
		return false
	}
	return strings.HasPrefix(sampleValue, "age-encryption-----")
}
