package cmd

import (
	"os"

	"github.com/rapour/cantainer/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "creates and runs a new container",
	Run: func(cmd *cobra.Command, args []string) {

		tempDir := internal.CreateTempDir()
		defer os.RemoveAll(tempDir)

		internal.Extract(tempDir)
		internal.Chroot(tempDir, "/bin/busybox", "/bin/ash")

	},
}
