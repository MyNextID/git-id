package cmd

import (
	"fmt"

	"github.com/mynextid/git-id/identity"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load [path]",
	Short: "Load an existing identity from the specified path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		id, err := identity.ReadIdentity(path)
		if err != nil {
			fmt.Printf("Error loading identity: %v\n", err)
			return
		}
		pubPEM, err := formatPublicKeyPEM(id.PublicKey)
		if err != nil {
			fmt.Printf("Error formatting public key: %v\n", err)
			return
		}
		fmt.Printf("Loaded identity.\n \n %s \n", pubPEM)
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
