package commands

import (
	"os"

	"github.com/rapour/cantainer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "creates and runs a new container",
	Run: func(cmd *cobra.Command, args []string) {

		tempDir := cantainer.CreateTempDir()
		defer os.RemoveAll(tempDir)

		cantainer.Extract(tempDir)
		cantainer.Chroot(tempDir, "/bin/busybox", "/bin/ash")

	},
}
