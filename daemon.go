package cantainer

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func RunDaemon(ctx context.Context, state *state) error {

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

	return g.Wait()
}
