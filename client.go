package cantainer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/netip"
	"time"
)

type Client struct {
	client  *http.Client
	Address string
}

// TODO: tune client timeouts
// TODO: add tracing to http client
func NewClient() *Client {
	return &Client{
		Address: "http://127.0.0.1:20043",
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (c *Client) RegisterContainer(ctx context.Context, network *netip.Prefix) (netip.Addr, error) {

	reqBody := RegisterContainerHTTPRequest{
		Network: *network,
	}

	var buf *bytes.Buffer
	err := json.NewEncoder(buf).Encode(reqBody)
	if err != nil {
		return netip.Addr{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Address, buf)
	if err != nil {
		return netip.Addr{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return netip.Addr{}, err
	}
	defer resp.Body.Close()

	var response RegisterContainerHTTPResponse
	if err := json.NewDecoder(req.Body).Decode(&resp); err != nil {
		return netip.Addr{}, err
	}

	return response.Address, nil
}
