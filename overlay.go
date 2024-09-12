package cantainer

import (
	"context"
	"log/slog"
	"time"
)

const (
	VXLAN_NAME  = "vxlan-host"
	VXLAN_ID    = 100
	BRIDGE_NAME = "bridge-host"
)

type overlayNetworkController struct {
	state                *state
	tickInterval         time.Duration
	CurrentPeerAddresses map[string]bool
}

func NewOverlayNetworkController(state *state) *overlayNetworkController {
	return &overlayNetworkController{
		state:                state,
		tickInterval:         500 * time.Millisecond,
		CurrentPeerAddresses: make(map[string]bool),
	}
}

func (onc *overlayNetworkController) Run(ctx context.Context) error {

	err := CreateBridge(BRIDGE_NAME)
	if err != nil {
		return err
	}

	err = CreateVXLan(VXLAN_NAME, VXLAN_ID, BRIDGE_NAME)
	if err != nil {
		return err
	}

	cli, err := onc.state.app.Client(ctx)
	if err != nil {
		return err
	}
	defer cli.Close()

	ticker := time.NewTicker(onc.tickInterval)

	for {
		select {
		case <-ticker.C:

			peers, err := onc.state.PeerAddresses()
			if err != nil {
				slog.Error("error returning peer addresses", slog.String("error", err.Error()))
				continue
			}

			for peer := range peers {
				if _, ok := onc.CurrentPeerAddresses[peer]; !ok {
					err := AddRemoteToVXLan(VXLAN_NAME, peer)
					if err != nil {
						slog.Error("error adding peer address", slog.String("error", err.Error()))
						continue
					}
					onc.CurrentPeerAddresses[peer] = true
				}
			}

			for currentPeer := range onc.CurrentPeerAddresses {
				if _, ok := peers[currentPeer]; !ok {
					err := RemoveFromVXLan(VXLAN_NAME, currentPeer)
					if err != nil {
						slog.Error("error removing peer address", slog.String("error", err.Error()))
						continue
					}

					delete(onc.CurrentPeerAddresses, currentPeer)
				}
			}

		case <-ctx.Done():
			return nil
		}
	}

}
