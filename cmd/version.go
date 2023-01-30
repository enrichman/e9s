package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: `E9s is a CLI to view and manage your Epinio cluster.`,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf("%s \nVersion: %s\n", AppName, Version)
		},
	}

	return versionCmd
}
