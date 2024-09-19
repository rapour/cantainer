package commands

import (
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
		defer cantainer.DeleteNetworkNamespace(cName)

		cantainer.Extract(tempDir)
		cantainer.NewContainer(cName, tempDir, "/bin/busybox", "/bin/ash")

	},
}
