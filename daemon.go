package cantainer

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type core struct {
	state *state
}

func NewCore(state *state) *core {
	return &core{state: state}
}

func (c *core) RunDaemon(ctx context.Context) error {

	onc := NewOverlayNetworkController(c.state)

	if err := c.state.RegisterNode(); err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)

	// graceful shutdown
	g.Go(func() error {
		select {

		case <-gCtx.Done():

			err := c.state.UnregisterNode()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			stateShutdownErr := c.state.Shutdown(shutdownCtx)
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
