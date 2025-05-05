package cmd

import (
	"fmt"

	"github.com/mynextid/git-id/identity"
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
		pubPEM, err := formatPublicKeyPEM(pubKey)
		if err != nil {
			fmt.Printf("Error formatting public key: %v\n", err)
			return
		}
		fmt.Printf("Loaded identity.\n \n %s \n", pubPEM)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
