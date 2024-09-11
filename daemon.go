package cantainer

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func RunDaemon(ctx context.Context, state *state) error {

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		select {

		case <-gCtx.Done():
			return nil
		}
	})

	return g.Wait()
}
