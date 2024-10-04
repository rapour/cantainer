package cantainer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"slices"
	"time"

	"golang.org/x/sync/errgroup"
)

type core struct {
	state *state
}

func NewCore(state *state) *core {
	return &core{state: state}
}

func (c *core) RegisterContainer(ctx context.Context, network *netip.Prefix) (netip.Addr, error) {

	containers, err := c.state.Containers()
	if err != nil {
		return netip.Addr{}, err
	}

	networkContainers := []netip.Addr{}
	for _, c := range containers {
		if network.Contains(c) {
			networkContainers = append(networkContainers, c)
		}
	}

	currentCandidate := network.Addr()
	for {
		if !slices.Contains(networkContainers, currentCandidate) {
			break
		}

		currentCandidate = currentCandidate.Next()
		if !currentCandidate.IsValid() {
			return currentCandidate, fmt.Errorf("network is full")
		}
	}

	return currentCandidate, c.state.RegisterContainer(&currentCandidate)
}

func (c *core) RunDaemon(ctx context.Context) error {

	onc := NewOverlayNetworkController(c.state)

	http := NewHTTP(c)

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
		return http.Run(gCtx)
	})

	g.Go(func() error {
		return onc.Run(gCtx)
	})

	return g.Wait()
}
