package cantainer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	url := fmt.Sprintf("%s/register", c.Address)

	reqBody := RegisterContainerHTTPRequest{
		Network: *network,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)
	if err != nil {
		return netip.Addr{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return netip.Addr{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return netip.Addr{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return netip.Addr{}, err
		}

		return netip.Addr{}, fmt.Errorf("daemon responsed with [code: %v]: %v", resp.StatusCode, string(body))
	}

	var response RegisterContainerHTTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return netip.Addr{}, err
	}

	return response.Address, nil
}
