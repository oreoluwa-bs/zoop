package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Deletes data in store",
	Long:  `Deletes data in store with a key(identifier)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if key == "" {
			fmt.Fprintf(os.Stderr, "Error: Key cannot be empty\n")
			os.Exit(1)
		}
		err := storeManager.Store.Delete(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	},
}
