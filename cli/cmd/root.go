package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "git-id",
	Short: "Git ID management CLI",
	Long:  `A CLI tool to generate, load, and fetch Ed25519 Git IDs.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
