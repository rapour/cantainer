package commands

import (
	"log/slog"
	"os"

	"github.com/rapour/cantainer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "creates and runs a new container",
	Run: func(cmd *cobra.Command, args []string) {

		tempDir, cName := cantainer.CreateTempDir()
		defer os.RemoveAll(tempDir)

		err := cantainer.CreateNetworkNamespace(cName)
		if err != nil {
			panic(err)
		}

		cantainer.Extract(tempDir)
		cantainer.Chroot(tempDir, "/bin/busybox", "/bin/ash")

		slog.Info("running cleanups")
		err = cantainer.DeleteNetworkNamespace(cName)
		if err != nil {
			panic(err)
		}

	},
}
