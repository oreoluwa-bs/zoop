package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Sets data in store",
	Long:  `Sets data in store with a key(identifier) and a value`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]
		if key == "" {
			fmt.Fprintf(os.Stderr, "Error: Key cannot be empty\n")
			os.Exit(1)
		}
		if value == "" {
			fmt.Fprintf(os.Stderr, "Error: Value cannot be empty\n")
			os.Exit(1)
		}

		err := storeManager.Store.Set(key, value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}
