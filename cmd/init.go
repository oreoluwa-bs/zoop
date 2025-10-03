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

var (
	initEncrypt bool
	initKeyFile string
	initForce   bool
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&initEncrypt, "encrypt", true, "Enable encryption after initialization")
	initCmd.Flags().StringVar(&initKeyFile, "key-file", "", "Path to key file (default: ~/.zoop/key.txt)")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Force regeneration of keys even if they exist")
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

		keyFilePath := cfg.KeyFile
		if initKeyFile != "" {
			keyFilePath = initKeyFile
		}

		var identity *age.X25519Identity

		// Check if key file exists and we can reuse it
		if !initForce {
			if keyData, err := os.ReadFile(keyFilePath); err == nil {
				if parsedIdentity, err := age.ParseX25519Identity(string(keyData)); err == nil {
					identity = parsedIdentity
					fmt.Printf("✅ Reusing existing keys from: %s\n", keyFilePath)
				}
			}
		}

		// Generate new keys if not reused
		if identity == nil {
			newIdentity, err := age.GenerateX25519Identity()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating key: %v\n", err)
				os.Exit(1)
			}
			identity = newIdentity
			if initForce {
				fmt.Printf("✅ Force regenerating keys: %s\n", keyFilePath)
			} else {
				fmt.Printf("✅ Generated new keys: %s\n", keyFilePath)
			}
		}

		if err := os.MkdirAll(filepath.Dir(keyFilePath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
			os.Exit(1)
		}

		keyFile, err := os.OpenFile(keyFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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
		configPath := filepath.Join(filepath.Dir(keyFilePath), "config.yaml")
		v.Set("encryption", initEncrypt)
		v.Set("key_file", keyFilePath)
		if err := v.WriteConfigAs(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config: %v\n", err)
			os.Exit(1)
		}

		if initEncrypt {
			fmt.Printf("✅ Encryption enabled! Key file: %s\n", keyFilePath)
		} else {
			fmt.Printf("✅ Initialized with encryption disabled. Key file: %s\n", keyFilePath)
			fmt.Println("💡 Enable encryption later with: zoop config set encryption true")
		}
	},
}
