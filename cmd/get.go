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
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if key == "" {
			fmt.Fprintf(os.Stderr, "Error: Key cannot be empty\n")
			os.Exit(1)
		}
		val, err := storeManager.Store.Get(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(val)
	},
}
