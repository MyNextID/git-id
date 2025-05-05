package cmd

import (
	"fmt"

	"github.com/mynextid/git-id/identity"
	"github.com/spf13/cobra"
)

var overwrite bool

var generateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "Generate a new identity at the specified path",
	Run: func(cmd *cobra.Command, args []string) {
		path := "secret-key.pem"
		if len(args) > 0 {
			path = args[0]

		}
		id, err := identity.GenerateIdentity(path, overwrite)
		if err != nil {
			fmt.Printf("Error generating identity: %v\n", err)
			return
		}
		fmt.Printf("Identity generated.\nPublic key: %x\n", id.PublicKey)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&overwrite, "force", "f", false, "Overwrite existing key files if they exist")

}
