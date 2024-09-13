package commands

import (
	"github.com/rapour/cantainer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(childCmd)
}

var childCmd = &cobra.Command{
	Use:   "child",
	Short: "internal command to spawn a child process",
	Run: func(cmd *cobra.Command, args []string) {

		cantainer.Child(args[0], "/bin/busybox", "/bin/ash")

	},
}
