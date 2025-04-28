package cmd

import (
	"fmt"

	"github.com/mynextid/gid/identity"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [GitHub handler]",
	Short: "Fetch a public key from a GitHub repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handler := args[0]

		ref := "gid/main"
		keyPath := "gid.pem"

		pubKey, err := identity.FetchPublicKeyFromGitHub(handler, ref, keyPath)
		if err != nil {
			fmt.Printf("Error fetching public key: %v\n", err)
			return
		}
		fmt.Printf("Fetched public key: %x\n", pubKey)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
