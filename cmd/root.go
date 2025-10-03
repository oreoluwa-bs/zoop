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
	store := internal.NewInMemoryStore()
	storeManager = internal.NewStoreManager(store)
}
