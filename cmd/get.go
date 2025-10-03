package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Gets data in store",
	Long:  `Gets data in store with a key(identifier)`,
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		status := needsMigration()
		if status == MigrationNeedsDecrypt {
			fmt.Println("⚠️  Detected encrypted data but encryption is disabled")
			fmt.Println("Run 'zoop migrate decrypt' to decrypt your data")
			fmt.Println("Or enable encryption with 'zoop config set encryption true'")
			os.Exit(1)
		} else if status == MigrationNeedsEncrypt {
			fmt.Println("⚠️  Detected plain data but encryption is enabled")
			fmt.Println("Run 'zoop migrate encrypt' to encrypt your data")
			fmt.Println("Or disable encryption with 'zoop config set encryption false'")
			os.Exit(1)
		}

		key := args[0]
		if key == "" {
			fmt.Fprintf(os.Stderr, "Error: Key cannot be empty\n")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		val, err := storeManager.Store.Get(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(val)
	},
}
