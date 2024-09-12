package cantainer

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

func RunDaemon(ctx context.Context, state *state) error {

	onc := NewOverlayNetworkController(state)

	if err := state.RegisterNode(); err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)

	// graceful shutdown
	g.Go(func() error {
		select {

		case <-gCtx.Done():

			err := state.UnregisterNode()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			stateShutdownErr := state.Shutdown(shutdownCtx)
			if stateShutdownErr != nil {
				err = errors.Join(err, stateShutdownErr)
			}

			if err == nil {
				slog.Info("gracefully shutdown")
			}

			return err
		}
	})

	g.Go(func() error {
		return onc.Run(gCtx)
	})

	return g.Wait()
}
