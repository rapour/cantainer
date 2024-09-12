package cantainer

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func RunDaemon(ctx context.Context, state *state) error {

	onc := NewOverlayNetworkController(state)

	if err := state.RegisterNode(); err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		select {

		case <-gCtx.Done():
			return nil
		}
	})

	g.Go(func() error {
		return onc.Run(gCtx)
	})

	return g.Wait()
}
