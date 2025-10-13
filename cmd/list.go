package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys in store",
	Long:  `List all keys in store`,
	Run: func(cmd *cobra.Command, args []string) {

		val, err := storeManager.Store.GetAllKeys()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for _, v := range val {
			fmt.Printf("%s\n", v)
		}
	},
}
