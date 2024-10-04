package commands

import (
	"context"
	"net"
	"net/netip"
	"os"

	"github.com/rapour/cantainer"
	"github.com/spf13/cobra"
)

var (
	network net.IPNet
)

func init() {
	rootCmd.AddCommand(newCmd)

	daemonCmd.Flags().IPNetVarP(&network, "network", "n", net.IPNet{}, "network CIDR to attach the container, if empty, container will not be able to reach other containers in the overlay network")

}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "creates and runs a new container",
	Run: func(cmd *cobra.Command, args []string) {

		tempDir, cName := cantainer.CreateTempDir()
		defer os.RemoveAll(tempDir)

		if err := cantainer.CreateNetworkNamespace(cName); err != nil {
			panic(err)
		}
		defer cantainer.DeleteNetworkNamespace(cName)

		if err := cantainer.ConnectNetworkNamespaceToBridge(cName, cantainer.BRIDGE_NAME); err != nil {
			panic(err)
		}

		if network.IP != nil {
			netIP, err := netip.ParsePrefix(network.String())
			if err != nil {
				panic(err)
			}

			address, err := cantainer.NewClient().RegisterContainer(context.Background(), &netIP)
			if err != nil {
				panic(err)
			}

			err = cantainer.AssignAddressToNamespace(cName, address)
			if err != nil {
				panic(err)
			}
		}

		cantainer.Extract(tempDir)
		cantainer.NewContainer(cName, tempDir, "/bin/busybox", "/bin/ash")

	},
}
