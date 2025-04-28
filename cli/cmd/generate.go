package cmd

import (
	"fmt"

	"github.com/mynextid/gid/identity"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "Generate a new identity at the specified path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		id, err := identity.GenerateIdentity(path)
		if err != nil {
			fmt.Printf("Error generating identity: %v\n", err)
			return
		}
		fmt.Printf("Identity generated.\nPublic key: %x\n", id.PublicKey)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
