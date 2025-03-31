package kvcli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ServerAddress is the default address for connecting to the KVD server
var ServerAddress string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kv",
	Short: "kv - a simple CLI to store key value pairs in the kvd service",
	Long: `kv is a command line interface for interacting with the kvd key-value store service.
It allows you to set, get, and delete keys, as well as view metrics about the store.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Display help if no subcommand is provided
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// init initializes configuration settings
func init() {
	// Define persistent flags used by all commands
	rootCmd.PersistentFlags().StringVar(&ServerAddress, "server", "http://localhost:8080", "Server address (e.g., http://localhost:8080)")
}
