package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Zoop",
	Long:  `All software has versions. This is Zoop's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Zoop v0.9 -- HEAD")
	},
}
