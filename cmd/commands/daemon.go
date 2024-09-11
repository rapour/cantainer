package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rapour/cantainer"
	"github.com/spf13/cobra"
)

var (
	seeds []string
	dir   string
)

func init() {
	rootCmd.AddCommand(daemonCmd)

	daemonCmd.Flags().StringSliceVarP(&seeds, "seeds", "s", nil, "seed servcer to connect (required)")
	daemonCmd.Flags().StringVarP(&dir, "dir", "d", "/tmp/cantainer-dqlite", "data directory")

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

		state, err := cantainer.NewState(addr, dir, seeds)
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		return cantainer.RunDaemon(ctx, state)
	},
}
