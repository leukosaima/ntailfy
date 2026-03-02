package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type TailscaleClient struct {
	oauthClientID     string
	oauthClientSecret string
	oauthScope        string
	tailnet           string
	client            *http.Client

	tokenMu           sync.Mutex
	oauthAccessToken  string
	oauthAccessExpiry time.Time
}

type Device struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Hostname           string `json:"hostname"`
	ConnectedToControl bool   `json:"connectedToControl"`
}

func (d *Device) Online() bool {
	return d.ConnectedToControl
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func NewTailscaleClient(config *Config) *TailscaleClient {
	return &TailscaleClient{
		oauthClientID:     config.TailscaleOAuthClientID,
		oauthClientSecret: config.TailscaleOAuthClientSecret,
		oauthScope:        config.TailscaleOAuthScope,
		tailnet:           config.TailscaleTailnet,
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

	bearer, err := t.oauthAccessTokenValue()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("Accept", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(b))
		if msg != "" {
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, msg)
		}
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var devicesResp DevicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&devicesResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return devicesResp.Devices, nil
}

func (t *TailscaleClient) oauthAccessTokenValue() (string, error) {
	t.tokenMu.Lock()
	defer t.tokenMu.Unlock()

	if t.oauthAccessToken != "" && time.Now().Add(1*time.Minute).Before(t.oauthAccessExpiry) {
		return t.oauthAccessToken, nil
	}

	tr, err := t.mintOAuthToken()
	if err != nil {
		return "", err
	}
	if tr.AccessToken == "" {
		return "", fmt.Errorf("tailscale oauth token response missing access_token")
	}

	t.oauthAccessToken = tr.AccessToken
	if tr.ExpiresIn > 0 {
		t.oauthAccessExpiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	} else {
		t.oauthAccessExpiry = time.Now().Add(55 * time.Minute)
	}

	return t.oauthAccessToken, nil
}

func (t *TailscaleClient) mintOAuthToken() (*oauthTokenResponse, error) {
	form := url.Values{}
	form.Set("client_id", t.oauthClientID)
	form.Set("client_secret", t.oauthClientSecret)
	if t.oauthScope != "" {
		form.Set("scope", t.oauthScope)
	}

	req, err := http.NewRequest("POST", "https://api.tailscale.com/api/v2/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating oauth token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting oauth token: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading oauth token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		msg := strings.TrimSpace(string(b))
		if msg != "" {
			return nil, fmt.Errorf("oauth token endpoint returned status %d: %s", resp.StatusCode, msg)
		}
		return nil, fmt.Errorf("oauth token endpoint returned status %d", resp.StatusCode)
	}

	var tr oauthTokenResponse
	if err := json.Unmarshal(b, &tr); err != nil {
		return nil, fmt.Errorf("decoding oauth token response: %w", err)
	}
	return &tr, nil
}
