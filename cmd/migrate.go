package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	internal "github.com/oreoluwa-bs/zoop/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

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

		if _, err := os.Stat(cfg.KeyFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Key file not found: %s\n", cfg.KeyFile)
			os.Exit(1)
		}

		keyData, err := os.ReadFile(cfg.KeyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading key file: %v\n", err)
			os.Exit(1)
		}

		cipher, err := internal.NewAgeCipher(string(keyData))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating cipher: %v\n", err)
			os.Exit(1)
		}

		plainStore, err := internal.NewJSONStore(cfg.DataFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading plain store: %v\n", err)
			os.Exit(1)
		}

		encryptedStore := internal.NewEncryptedStore(plainStore, cipher)

		if action == "decrypt" {
			err := internal.MigrateStores(encryptedStore, plainStore)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error migrating: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("✅ Data decrypted successfully")

		} else if action == "encrypt" {
			err := internal.MigrateStores(plainStore, encryptedStore)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error migrating: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("✅ Data encrypted successfully")
		}
	},
}

type MigrationStatus int

const (
	MigrationNone MigrationStatus = iota
	MigrationNeedsEncrypt
	MigrationNeedsDecrypt
)

func needsMigration() MigrationStatus {
	cfg, err := internal.LoadConfig()
	if err != nil {
		return MigrationNone
	}
	data, err := os.ReadFile(cfg.DataFile)
	if err != nil {
		return MigrationNone
	}
	var src internal.JSONDatasource
	if err := json.Unmarshal(data, &src); err != nil {
		return MigrationNone // corrupted
	}

	hasEncrypted := false
	hasPlain := false
	for _, v := range src.Data {
		if _, err := base64.StdEncoding.DecodeString(v); err == nil && len(v) > 10 {
			hasEncrypted = true
		} else {
			hasPlain = true
		}
	}

	if hasEncrypted && !hasPlain {
		// All encrypted
		if !cfg.Encryption {
			return MigrationNeedsDecrypt
		}
	} else if hasPlain && !hasEncrypted {
		// All plain
		if cfg.Encryption {
			return MigrationNeedsEncrypt
		}
	} else if hasEncrypted && hasPlain {
		// Mixed - need to handle, but for now, if encryption enabled, assume needs encrypt for plain, but complicated
		// For simplicity, if cfg.Encryption, return NeedsEncrypt (to encrypt remaining plain)
		// Else, return NeedsDecrypt (to decrypt remaining encrypted)
		if cfg.Encryption {
			return MigrationNeedsEncrypt
		} else {
			return MigrationNeedsDecrypt
		}
	}

	return MigrationNone
}
