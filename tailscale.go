package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TailscaleClient struct {
	apiKey  string
	tailnet string
	client  *http.Client
}

type ClientConnectivity struct {
	ConnectedToControl bool `json:"connectedToControl"`
}

type Device struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Hostname           string              `json:"hostname"`
	ClientConnectivity *ClientConnectivity `json:"clientConnectivity"`
}

// Online returns true if the device is connected to the control plane
func (d *Device) Online() bool {
	if d.ClientConnectivity == nil {
		return false
	}
	return d.ClientConnectivity.ConnectedToControl
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

func NewTailscaleClient(apiKey, tailnet string) *TailscaleClient {
	return &TailscaleClient{
		apiKey:  apiKey,
		tailnet: tailnet,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (t *TailscaleClient) GetDevices() ([]Device, error) {
	url := fmt.Sprintf("https://api.tailscale.com/api/v2/tailnet/%s/devices", t.tailnet)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var devicesResp DevicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&devicesResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return devicesResp.Devices, nil
}
