package commands

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rapour/cantainer"
	"github.com/spf13/cobra"
)

var (
	seeds  []string
	dir    string
	socket cantainer.Socket
)

func init() {
	rootCmd.AddCommand(daemonCmd)

	daemonCmd.Flags().StringSliceVarP(&seeds, "seeds", "s", nil, "seed server to connect (required)")
	daemonCmd.Flags().StringVarP(&dir, "dir", "d", "/tmp/cantainer-dqlite", "data directory")
	daemonCmd.Flags().UintVarP(&socket.Port, "port", "p", 9000, "port number that the daemon will listen on")

	daemonCmd.MarkFlagRequired("seeds")
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "runs the cantainer daemon",
	RunE: func(cmd *cobra.Command, args []string) error {

		addr, err := cantainer.GetAddress()
		if err != nil {
			return err
		}
		slog.Info("obtained host address", slog.Any("address", addr))
		socket.Address = addr

		state, err := cantainer.NewState(socket, dir, seeds)
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		core := cantainer.NewCore(state)

		return core.RunDaemon(ctx)
	},
}
